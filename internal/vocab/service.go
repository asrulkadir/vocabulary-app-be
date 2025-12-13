package vocab

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrVocabNotFound     = errors.New("vocabulary not found")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrNoVocabsAvailable = errors.New("no vocabularies available for testing")
)

// Service handles business logic for vocabulary
type Service interface {
	Create(ctx context.Context, userID string, req *CreateVocabRequest) (*Vocabulary, error)
	GetByID(ctx context.Context, userID, id string) (*Vocabulary, error)
	GetByUserID(ctx context.Context, userID string, page, pageSize int, search, status string) (*VocabListResponse, error)
	Update(ctx context.Context, userID, id string, req *UpdateVocabRequest) (*Vocabulary, error)
	Delete(ctx context.Context, userID, id string) error
	GetRandomForTest(ctx context.Context, userID string, status string) (*TestVocabulary, error)
	GetTestOptions(ctx context.Context, userID string, vocabID string) (*TestOptionsResponse, error)
	GetVocabStats(ctx context.Context, userID string) (map[string]int64, error)
	ValidateTestAnswer(ctx context.Context, userID, id string, input string) (*TestResultResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new vocabulary service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create creates a new vocabulary entry
func (s *service) Create(ctx context.Context, userID string, req *CreateVocabRequest) (*Vocabulary, error) {
	vocab := &Vocabulary{
		UserID:      userID,
		Word:        req.Word,
		Definition:  req.Definition,
		Example:     req.Example,
		Translation: req.Translation,
		Status:      StatusLearning,
		TestCount:   0,
		PassedTestCount: 0,
		FailedTestCount: 0,
	}

	if err := s.repo.Create(ctx, vocab); err != nil {
		return nil, err
	}

	return vocab, nil
}

// GetByID retrieves a vocabulary by ID
func (s *service) GetByID(ctx context.Context, userID, id string) (*Vocabulary, error) {
	vocab, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if vocab == nil {
		return nil, ErrVocabNotFound
	}

	// Check ownership
	if vocab.UserID != userID {
		return nil, ErrUnauthorized
	}

	return vocab, nil
}

// GetByUserID retrieves vocabularies by user ID with pagination
func (s *service) GetByUserID(ctx context.Context, userID string, page, pageSize int, search, status string) (*VocabListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	vocabularies, total, err := s.repo.FindByUserID(ctx, userID, page, pageSize, search, status)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &VocabListResponse{
		Data:       vocabularies,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Search:     search,
		Status:     status,
	}, nil
}

// Update updates a vocabulary entry
func (s *service) Update(ctx context.Context, userID, id string, req *UpdateVocabRequest) (*Vocabulary, error) {
	vocab, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if vocab == nil {
		return nil, ErrVocabNotFound
	}

	// Check ownership
	if vocab.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.Word != "" {
		vocab.Word = req.Word
	}
	exampleValue, err := req.Example.Value()
	if err == nil && exampleValue != "" {
		vocab.Example = req.Example
	}

	vocab.Definition = req.Definition
	vocab.Translation = req.Translation

	if err := s.repo.Update(ctx, vocab); err != nil {
		return nil, err
	}

	return vocab, nil
}

// Delete deletes a vocabulary entry
func (s *service) Delete(ctx context.Context, userID, id string) error {
	vocab, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if vocab == nil {
		return ErrVocabNotFound
	}

	// Check ownership
	if vocab.UserID != userID {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, id)
}

// GetRandomForTest gets a random vocabulary for testing with optional status filter
func (s *service) GetRandomForTest(ctx context.Context, userID string, status string) (*TestVocabulary, error) {
	vocab, err := s.repo.FindRandomByUserIDAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}
	if vocab == nil {
		return nil, ErrNoVocabsAvailable
	}
	return vocab.ToTestVocabulary(), nil
}

// GetTestOptions gets random vocabulary options for multiple-choice test (4 total: 1 correct + 3 wrong)
func (s *service) GetTestOptions(ctx context.Context, userID string, vocabID string) (*TestOptionsResponse, error) {
	// Get the correct answer (the vocabulary being tested)
	correctVocab, err := s.repo.FindByID(ctx, vocabID)
	if err != nil {
		return nil, err
	}
	if correctVocab == nil {
		return nil, ErrVocabNotFound
	}

	// Check ownership
	if correctVocab.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Get 3 random wrong options excluding the correct answer
	wrongOptions, err := s.repo.FindRandomOptionsExcluding(ctx, userID, vocabID, 3)
	if err != nil {
		return nil, err
	}

	// Create TestOption array with correct answer + wrong answers
	testOptions := make([]TestOption, 0, 4)
	
	// Add correct answer
	testOptions = append(testOptions, TestOption{
		ID:          correctVocab.ID,
		Translation: correctVocab.Translation,
	})
	
	// Add wrong answers
	for _, opt := range wrongOptions {
		testOptions = append(testOptions, TestOption{
			ID:          opt.ID,
			Translation: opt.Translation,
		})
	}

	return &TestOptionsResponse{
		Options: testOptions,
	}, nil
}

// GetVocabStats gets vocabulary statistics for the user
func (s *service) GetVocabStats(ctx context.Context, userID string) (map[string]int64, error) {
	stats := make(map[string]int64)

	total, err := s.repo.CountByUserIDAndStatus(ctx, userID, "")
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	learning, err := s.repo.CountByUserIDAndStatus(ctx, userID, string(StatusLearning))
	if err != nil {
		return nil, err
	}
	stats["learning"] = learning

	memorized, err := s.repo.CountByUserIDAndStatus(ctx, userID, string(StatusMemorized))
	if err != nil {
		return nil, err
	}
	stats["memorized"] = memorized

	return stats, nil
}

// ValidateTestAnswer validates the user's answer and updates test result
func (s *service) ValidateTestAnswer(ctx context.Context, userID, id string, input string) (*TestResultResponse, error) {
	vocab, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if vocab == nil {
		return nil, ErrVocabNotFound
	}

	// Check ownership
	if vocab.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Validate answer - compare with translation (case-insensitive, trimmed)
	correctAnswer := vocab.Translation
	passed := strings.EqualFold(strings.TrimSpace(input), strings.TrimSpace(correctAnswer))

	// Update test counts
	vocab.TestCount++
	if passed {
		vocab.PassedTestCount++
	} else {
		vocab.FailedTestCount++
	}

	// Auto-memorize if passed - failed >= 10
	if vocab.PassedTestCount-vocab.FailedTestCount >= 10 && vocab.Status != StatusMemorized {
		vocab.Status = StatusMemorized
	} else if vocab.PassedTestCount-vocab.FailedTestCount < 10 && vocab.Status != StatusLearning {
		vocab.Status = StatusLearning
	}

	if err := s.repo.Update(ctx, vocab); err != nil {
		return nil, err
	}

	response := &TestResultResponse{
		Passed:     passed,
		Vocabulary: *vocab.ToTestVocabulary(),
	}

	// Include correct answer if failed
	if !passed {
		response.CorrectAnswer = correctAnswer
	}

	return response, nil
}
