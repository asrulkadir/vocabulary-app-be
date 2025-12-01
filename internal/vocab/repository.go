package vocab

import (
	"context"
	"database/sql"
)

// Repository handles data access for vocabulary
type Repository interface {
	Create(ctx context.Context, vocab *Vocabulary) error
	FindByID(ctx context.Context, id int64) (*Vocabulary, error)
	FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]Vocabulary, int64, error)
	Update(ctx context.Context, vocab *Vocabulary) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new vocabulary repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Create creates a new vocabulary entry
func (r *repository) Create(ctx context.Context, vocab *Vocabulary) error {
	query := `INSERT INTO vocabularies (user_id, word, definition, example, translation, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		vocab.UserID,
		vocab.Word,
		vocab.Definition,
		vocab.Example,
		vocab.Translation,
	).Scan(&vocab.ID)
}

// FindByID finds a vocabulary by ID
func (r *repository) FindByID(ctx context.Context, id int64) (*Vocabulary, error) {
	query := `SELECT id, user_id, word, definition, example, translation, created_at, updated_at 
			  FROM vocabularies WHERE id = $1`

	var vocab Vocabulary
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vocab.ID,
		&vocab.UserID,
		&vocab.Word,
		&vocab.Definition,
		&vocab.Example,
		&vocab.Translation,
		&vocab.CreatedAt,
		&vocab.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &vocab, nil
}

// FindByUserID finds vocabularies by user ID with pagination
func (r *repository) FindByUserID(ctx context.Context, userID int64, page, pageSize int) ([]Vocabulary, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM vocabularies WHERE user_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	query := `SELECT id, user_id, word, definition, example, translation, created_at, updated_at 
			  FROM vocabularies WHERE user_id = $1 
			  ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vocabularies []Vocabulary
	for rows.Next() {
		var vocab Vocabulary
		if err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Word,
			&vocab.Definition,
			&vocab.Example,
			&vocab.Translation,
			&vocab.CreatedAt,
			&vocab.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, total, nil
}

// Update updates a vocabulary entry
func (r *repository) Update(ctx context.Context, vocab *Vocabulary) error {
	query := `UPDATE vocabularies 
			  SET word = $1, definition = $2, example = $3, translation = $4, updated_at = NOW() 
			  WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		vocab.Word,
		vocab.Definition,
		vocab.Example,
		vocab.Translation,
		vocab.ID,
	)
	return err
}

// Delete deletes a vocabulary entry
func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM vocabularies WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
