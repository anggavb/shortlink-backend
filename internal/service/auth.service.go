package service

import (
	"context"

	"github.com/shortlink-backend/internal/dto"
	"github.com/shortlink-backend/internal/repository"
	"github.com/shortlink-backend/pkg"
)

type AuthService struct {
	authRepository *repository.AuthRepository
	authCache      *repository.AuthCacheRepository
	// mailer         Mailer
}

func NewAuthService(authRepository *repository.AuthRepository, authCache *repository.AuthCacheRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		authCache:      authCache,
		// mailer:         mailer,
	}
}

func (as *AuthService) RegisterUser(ctx context.Context, user dto.RegisterRequest) error {
	// hashing password
	var hash pkg.HashConfig
	hash.UseRecommended()

	hashedPassword := hash.GenerateHash(user.Password)
	if err := as.authRepository.AddNewUser(ctx, user.Email, hashedPassword); err != nil {
		return err
	}

	return nil
}
