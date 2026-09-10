package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/german/backend/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) GetByUsername(ctx context.Context, username string) (models.User, error) {
	return repository.get(ctx, `
		SELECT id, username, password_hash, role, expert_group_id,
		       full_name, avg_rating, max_tickets
		FROM users
		WHERE username = $1`, username)
}

func (repository *UserRepository) GetByID(ctx context.Context, id int64) (models.User, error) {
	return repository.get(ctx, `
		SELECT id, username, password_hash, role, expert_group_id,
		       full_name, avg_rating, max_tickets
		FROM users
		WHERE id = $1`, id)
}

func (repository *UserRepository) get(ctx context.Context, query string, argument any) (models.User, error) {
	var user models.User
	err := repository.db.QueryRowContext(ctx, query, argument).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.ExpertGroupID,
		&user.FullName,
		&user.AvgRating,
		&user.MaxTickets,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
