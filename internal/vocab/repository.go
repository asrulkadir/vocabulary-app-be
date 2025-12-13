package vocab

import (
	"context"
	"database/sql"
	"strconv"
)

// Repository handles data access for vocabulary
type Repository interface {
	Create(ctx context.Context, vocab *Vocabulary) error
	FindByID(ctx context.Context, id string) (*Vocabulary, error)
	FindByUserID(ctx context.Context, userID string, page, pageSize int, search, status string) ([]Vocabulary, int64, error)
	FindRandomByUserIDAndStatus(ctx context.Context, userID string, status string) (*Vocabulary, error)
	FindRandomOptionsExcluding(ctx context.Context, userID string, excludeID string, count int) ([]Vocabulary, error)
	CountByUserIDAndStatus(ctx context.Context, userID string, status string) (int64, error)
	Update(ctx context.Context, vocab *Vocabulary) error
	Delete(ctx context.Context, id string) error
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
	query := `INSERT INTO vocabularies (user_id, word, definition, example, translation, status, test_count, passed_test_count, failed_test_count, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		vocab.UserID,
		vocab.Word,
		vocab.Definition,
		vocab.Example,
		vocab.Translation,
		vocab.Status,
		vocab.TestCount,
		vocab.PassedTestCount,
		vocab.FailedTestCount,
	).Scan(&vocab.ID)
}

// FindByID finds a vocabulary by ID
func (r *repository) FindByID(ctx context.Context, id string) (*Vocabulary, error) {
	query := `SELECT id, user_id, word, definition, example, translation, status, test_count, passed_test_count, failed_test_count, created_at, updated_at 
			  FROM vocabularies WHERE id = $1`

	var vocab Vocabulary
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vocab.ID,
		&vocab.UserID,
		&vocab.Word,
		&vocab.Definition,
		&vocab.Example,
		&vocab.Translation,
		&vocab.Status,
		&vocab.TestCount,
		&vocab.PassedTestCount,
		&vocab.FailedTestCount,
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

// FindByUserID finds vocabularies by user ID with pagination, search, and status filter
func (r *repository) FindByUserID(ctx context.Context, userID string, page, pageSize int, search, status string) ([]Vocabulary, int64, error) {
	// Build dynamic query conditions
	baseCondition := "user_id = $1"
	args := []any{userID}
	argIndex := 2

	if search != "" {
		baseCondition += " AND (word ILIKE $" + itoa(argIndex) + " OR translation ILIKE $" + itoa(argIndex) + " OR definition ILIKE $" + itoa(argIndex) + ")"
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if status != "" && status != "all" {
		baseCondition += " AND status = $" + itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM vocabularies WHERE " + baseCondition
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	query := `SELECT id, user_id, word, definition, example, translation, status, test_count, passed_test_count, failed_test_count, created_at, updated_at 
			  FROM vocabularies WHERE ` + baseCondition + ` 
			  ORDER BY created_at DESC LIMIT $` + itoa(argIndex) + ` OFFSET $` + itoa(argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
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
			&vocab.Status,
			&vocab.TestCount,
			&vocab.PassedTestCount,
			&vocab.FailedTestCount,
			&vocab.CreatedAt,
			&vocab.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, total, nil
}

// itoa converts int to string for query building
func itoa(i int) string {
	return strconv.Itoa(i)
}

// Update updates a vocabulary entry
func (r *repository) Update(ctx context.Context, vocab *Vocabulary) error {
	query := `UPDATE vocabularies 
			  SET word = $1, definition = $2, example = $3, translation = $4, status = $5, test_count = $6, passed_test_count = $7, failed_test_count = $8, updated_at = NOW() 
			  WHERE id = $9`

	_, err := r.db.ExecContext(ctx, query,
		vocab.Word,
		vocab.Definition,
		vocab.Example,
		vocab.Translation,
		vocab.Status,
		vocab.TestCount,
		vocab.PassedTestCount,
		vocab.FailedTestCount,
		vocab.ID,
	)
	return err
}

// Delete deletes a vocabulary entry
func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM vocabularies WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// FindRandomByUserIDAndStatus finds a random vocabulary by user ID and optional status filter
func (r *repository) FindRandomByUserIDAndStatus(ctx context.Context, userID string, status string) (*Vocabulary, error) {
	var query string
	var args []any

	if status == "" || status == "all" {
		query = `SELECT id, user_id, word, definition, example, translation, status, test_count, passed_test_count, failed_test_count, created_at, updated_at 
				 FROM vocabularies WHERE user_id = $1 
				 ORDER BY RANDOM() LIMIT 1`
		args = []any{userID}
	} else {
		query = `SELECT id, user_id, word, definition, example, translation, status, test_count, passed_test_count, failed_test_count, created_at, updated_at 
				 FROM vocabularies WHERE user_id = $1 AND status = $2 
				 ORDER BY RANDOM() LIMIT 1`
		args = []any{userID, status}
	}

	var vocab Vocabulary
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&vocab.ID,
		&vocab.UserID,
		&vocab.Word,
		&vocab.Definition,
		&vocab.Example,
		&vocab.Translation,
		&vocab.Status,
		&vocab.TestCount,
		&vocab.PassedTestCount,
		&vocab.FailedTestCount,
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

// FindRandomOptionsExcluding finds random vocabularies excluding a specific ID (for multiple choice options)
func (r *repository) FindRandomOptionsExcluding(ctx context.Context, userID string, excludeID string, count int) ([]Vocabulary, error) {
	query := `SELECT id, user_id, word, definition, example, translation, status, test_count, passed_test_count, failed_test_count, created_at, updated_at 
			  FROM vocabularies WHERE user_id = $1 AND id != $2 AND translation IS NOT NULL AND translation != ''
			  ORDER BY RANDOM() LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, userID, excludeID, count)
	if err != nil {
		return nil, err
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
			&vocab.Status,
			&vocab.TestCount,
			&vocab.PassedTestCount,
			&vocab.FailedTestCount,
			&vocab.CreatedAt,
			&vocab.UpdatedAt,
		); err != nil {
			return nil, err
		}
		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

// CountByUserIDAndStatus counts vocabularies by user ID and optional status filter
func (r *repository) CountByUserIDAndStatus(ctx context.Context, userID string, status string) (int64, error) {
	var query string
	var args []any

	if status == "" || status == "all" {
		query = `SELECT COUNT(*) FROM vocabularies WHERE user_id = $1`
		args = []any{userID}
	} else {
		query = `SELECT COUNT(*) FROM vocabularies WHERE user_id = $1 AND status = $2`
		args = []any{userID, status}
	}

	var count int64
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
