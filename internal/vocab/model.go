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
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	Word              string    `json:"word"`
	Definition        string    `json:"definition"`
	Example           Examples  `json:"example,omitempty"`
	Translation       string    `json:"translation,omitempty"`
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

// TestResultRequest represents the test result update request
type TestResultRequest struct {
	Passed bool `json:"passed"`
}

// VocabListResponse represents the vocabulary list response
type VocabListResponse struct {
	Data       []Vocabulary `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
