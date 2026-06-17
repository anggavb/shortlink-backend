package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
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

// func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
// 	sql := `
// 		SELECT id, fullname, email, password, phone_number, photo, is_verified, pin
// 		FROM users
// 		WHERE email = $1;
// 	`
// 	args := []any{email}

// 	var user model.User
// 	if err := ar.db.QueryRow(ctx, sql, args...).Scan(
// 		&user.Id,
// 		&user.Fullname,
// 		&user.Email,
// 		&user.Password,
// 		&user.PhoneNumber,
// 		&user.Photo,
// 		&user.IsVerified,
// 		&user.Pin,
// 	); err != nil {
// 		return model.User{}, err
// 	}

// 	return user, nil
// }

// func (ar *AuthRepository) GetPasswordResetUserByEmail(ctx context.Context, email string) (model.User, error) {
// 	sql := `
// 		SELECT id, email
// 		FROM users
// 		WHERE email = $1;
// 	`
// 	args := []any{email}

// 	var user model.User
// 	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email); err != nil {
// 		return model.User{}, err
// 	}

// 	return user, nil
// }

// func (ar *AuthRepository) UpdatePasswordById(ctx context.Context, userId int, hashedPassword string) error {
// 	tx, err := ar.db.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback(ctx)

// 	userSQL := `
// 		UPDATE users
// 		SET
// 			password = $2,
// 			updated_at = NOW()
// 		WHERE id = $1;
// 	`
// 	userCmd, err := tx.Exec(ctx, userSQL, userId, hashedPassword)
// 	if err != nil {
// 		return err
// 	}
// 	if userCmd.RowsAffected() == 0 {
// 		return pgx.ErrNoRows
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		return err
// 	}

// 	return nil
// }
