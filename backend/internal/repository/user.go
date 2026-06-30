package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, login, password_hash, created_at)
		 VALUES ($1, $2, $3, $4)`,
		user.ID, user.Login, user.PasswordHash, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, login, password_hash, created_at FROM users WHERE login = $1`,
		login,
	).Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by login: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, login, password_hash, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}
