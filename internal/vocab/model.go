package vocab

import "time"

// Vocabulary represents the vocabulary domain model
type Vocabulary struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Word        string    `json:"word"`
	Definition  string    `json:"definition"`
	Example     string    `json:"example,omitempty"`
	Translation string    `json:"translation,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateVocabRequest represents the create vocabulary request payload
type CreateVocabRequest struct {
	Word        string `json:"word" binding:"required"`
	Definition  string `json:"definition" binding:"required"`
	Example     string `json:"example"`
	Translation string `json:"translation"`
}

// UpdateVocabRequest represents the update vocabulary request payload
type UpdateVocabRequest struct {
	Word        string `json:"word"`
	Definition  string `json:"definition"`
	Example     string `json:"example"`
	Translation string `json:"translation"`
}

// VocabListResponse represents the vocabulary list response
type VocabListResponse struct {
	Data       []Vocabulary `json:"data"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
