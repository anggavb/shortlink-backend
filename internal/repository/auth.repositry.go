package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shortlink-backend/internal/model"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) AddNewUser(ctx context.Context, email, hashedPassword string) error {
	sql := `
		INSERT INTO users
		(email, password)
		VALUES
		($1, $2);
	`
	args := []any{email, hashedPassword}

	if _, err := ar.db.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT id, email, password, first_name, last_name, role, workplace, is_member, is_receive_email, photo, verified_at, created_at, updated_at
		FROM users
		WHERE email = $1;
	`
	args := []any{email}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Workplace,
		&user.IsMember,
		&user.IsReceiveEmail,
		&user.Photo,
		&user.VerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return model.User{}, err
	}

	return user, nil
}
