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

func (as *AuthService) LoginUser(ctx context.Context, user dto.LoginRequest) (dto.AuthResponse, error) {
	userLogin, err := as.authRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	var hash pkg.HashConfig
	if err := hash.Compare(user.Password, userLogin.Password); err != nil {
		return dto.AuthResponse{}, err
	}

	claims := pkg.NewClaims(userLogin.Id, user.Email)
	token, err := claims.GenerateJWT()
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if err := as.authCache.SaveToken(ctx, userLogin, claims.ExpiresAt.Time); err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Token:     token,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (as *AuthService) LogoutUser(ctx context.Context, userId int) error {
	return as.authCache.DeleteToken(ctx, userId)
}
