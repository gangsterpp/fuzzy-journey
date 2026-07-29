package auth

import (
	"context"
	"database/sql"

	"github.com/gangsterpp/fuzzy-journey/internal/response"
	user "github.com/gangsterpp/fuzzy-journey/internal/user"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	var u user.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email,passwordHash
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) DeleteByID(
	ctx context.Context,
	id string,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return response.ErrUserNotFound
	}

	return nil
}

func (r *Repository) RegisterUser(
	ctx context.Context,
	email string,
	passwordHash string,
) (*user.User, error) {
	u := user.User{Email: email}
	query := `
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, created_at`
	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
		passwordHash,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func NewAuthRepository(db *sql.DB) Repository {
	return Repository{db: db}
}
