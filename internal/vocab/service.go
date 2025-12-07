package vocab

import (
	"context"
	"errors"
)

var (
	ErrVocabNotFound = errors.New("vocabulary not found")
	ErrUnauthorized  = errors.New("unauthorized access")
)

// Service handles business logic for vocabulary
type Service interface {
	Create(ctx context.Context, userID int64, req *CreateVocabRequest) (*Vocabulary, error)
	GetByID(ctx context.Context, userID, id int64) (*Vocabulary, error)
	GetByUserID(ctx context.Context, userID int64, page, pageSize int) (*VocabListResponse, error)
	Update(ctx context.Context, userID, id int64, req *UpdateVocabRequest) (*Vocabulary, error)
	Delete(ctx context.Context, userID, id int64) error
	UpdateTestResult(ctx context.Context, userID, id int64, passed bool) (*Vocabulary, error)
}

type service struct {
	repo Repository
}

// NewService creates a new vocabulary service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create creates a new vocabulary entry
func (s *service) Create(ctx context.Context, userID int64, req *CreateVocabRequest) (*Vocabulary, error) {
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
func (s *service) GetByID(ctx context.Context, userID, id int64) (*Vocabulary, error) {
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
func (s *service) GetByUserID(ctx context.Context, userID int64, page, pageSize int) (*VocabListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	vocabularies, total, err := s.repo.FindByUserID(ctx, userID, page, pageSize)
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
	}, nil
}

// Update updates a vocabulary entry
func (s *service) Update(ctx context.Context, userID, id int64, req *UpdateVocabRequest) (*Vocabulary, error) {
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
	if req.Definition != "" {
		vocab.Definition = req.Definition
	}
	exampleValue, err := req.Example.Value()
	if err == nil && exampleValue != "" {
		vocab.Example = req.Example
	}
	if req.Translation != "" {
		vocab.Translation = req.Translation
	}

	if err := s.repo.Update(ctx, vocab); err != nil {
		return nil, err
	}

	return vocab, nil
}

// Delete deletes a vocabulary entry
func (s *service) Delete(ctx context.Context, userID, id int64) error {
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

// UpdateTestResult updates the test result for a vocabulary
func (s *service) UpdateTestResult(ctx context.Context, userID, id int64, passed bool) (*Vocabulary, error) {
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
	}

	if err := s.repo.Update(ctx, vocab); err != nil {
		return nil, err
	}

	return vocab, nil
}
