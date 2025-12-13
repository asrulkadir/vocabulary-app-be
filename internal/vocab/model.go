package vocab

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Examples is a custom type for storing multiple examples as JSON
type Examples []string

// Scan implements the sql.Scanner interface
func (e *Examples) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &e)
}

// Value implements the driver.Valuer interface
func (e Examples) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Status represents the vocabulary learning status
type Status string

const (
	StatusLearning   Status = "learning"
	StatusMemorized  Status = "memorized"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	return s == StatusLearning || s == StatusMemorized
}

// String returns the string representation of status
func (s Status) String() string {
	return string(s)
}

// Vocabulary represents the vocabulary domain model
type Vocabulary struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Word              string    `json:"word"`
	Definition        string    `json:"definition"`
	Example           Examples  `json:"example,omitempty"`
	Translation       string    `json:"translation"`
	Status            Status    `json:"status"`
	TestCount         int64     `json:"test_count"`
	PassedTestCount   int64     `json:"passed_test_count"`
	FailedTestCount   int64     `json:"failed_test_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CreateVocabRequest represents the create vocabulary request payload
type CreateVocabRequest struct {
	Word        string   `json:"word" binding:"required"`
	Definition  string   `json:"definition"`
	Example     Examples `json:"example"`
	Translation string   `json:"translation"`
}

// UpdateVocabRequest represents the update vocabulary request payload
type UpdateVocabRequest struct {
	Word        string   `json:"word"`
	Definition  string   `json:"definition"`
	Example     Examples `json:"example"`
	Translation string   `json:"translation"`
	Status      Status   `json:"status"`
}

// TestResultRequest represents the test result update request (input-based validation)
type TestResultRequest struct {
	Input string `json:"input" binding:"required"`
}

// TestVocabulary represents vocabulary for testing (without answers)
type TestVocabulary struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Word            string    `json:"word"`
	Status          Status    `json:"status"`
	TestCount       int64     `json:"test_count"`
	PassedTestCount int64     `json:"passed_test_count"`
	FailedTestCount int64     `json:"failed_test_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ToTestVocabulary converts Vocabulary to TestVocabulary (hides answers)
func (v *Vocabulary) ToTestVocabulary() *TestVocabulary {
	return &TestVocabulary{
		ID:              v.ID,
		UserID:          v.UserID,
		Word:            v.Word,
		Status:          v.Status,
		TestCount:       v.TestCount,
		PassedTestCount: v.PassedTestCount,
		FailedTestCount: v.FailedTestCount,
		CreatedAt:       v.CreatedAt,
		UpdatedAt:       v.UpdatedAt,
	}
}

// TestOption represents a single option for multiple choice tests (only translations)
type TestOption struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
}

// TestOptionsResponse represents multiple choice options response
type TestOptionsResponse struct {
	Options []TestOption `json:"options"`
}

// TestResultResponse represents the test result response
type TestResultResponse struct {
	Passed        bool          `json:"passed"`
	CorrectAnswer string        `json:"correct_answer,omitempty"`
	Vocabulary    TestVocabulary `json:"vocabulary"`
}

// VocabListResponse represents the vocabulary list response
type VocabListResponse struct {
	Data       []Vocabulary `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
	Search     string       `json:"search,omitempty"`
	Status     string       `json:"status,omitempty"`
}
