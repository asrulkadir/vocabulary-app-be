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

// Vocabulary represents the vocabulary domain model
type Vocabulary struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Word        string    `json:"word"`
	Definition  string    `json:"definition"`
	Example     Examples  `json:"example,omitempty"`
	Translation string    `json:"translation,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateVocabRequest represents the create vocabulary request payload
type CreateVocabRequest struct {
	Word        string   `json:"word" binding:"required"`
	Definition  string   `json:"definition" binding:"required"`
	Example     Examples `json:"example"`
	Translation string   `json:"translation"`
}

// UpdateVocabRequest represents the update vocabulary request payload
type UpdateVocabRequest struct {
	Word        string   `json:"word"`
	Definition  string   `json:"definition"`
	Example     Examples `json:"example"`
	Translation string   `json:"translation"`
	Status      string   `json:"status"`
}

// VocabListResponse represents the vocabulary list response
type VocabListResponse struct {
	Data       []Vocabulary `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
