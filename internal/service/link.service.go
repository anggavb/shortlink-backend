package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/shortlink-backend/internal/dto"
	"github.com/shortlink-backend/internal/repository"
)

const (
	linkCacheTTL              = time.Hour
	maxSlugGenerationAttempts = 5
	randomSlugLength          = 5
)

var (
	ErrDuplicateSlug        = errors.New("slug already used")
	ErrGenerateSlugFailed   = errors.New("failed to generate unique slug")
	ErrDuplicateOriginalURL = errors.New("original URL already used")
)

type LinkService struct {
	linkRepository *repository.LinkRepository
	linkCache      *repository.LinkCacheRepository
}

func NewLinkService(linkRepository *repository.LinkRepository, linkCache *repository.LinkCacheRepository) *LinkService {
	return &LinkService{
		linkRepository: linkRepository,
		linkCache:      linkCache,
	}
}

func (ls *LinkService) CreateLink(ctx context.Context, userId int, body dto.CreateLinkRequest) (dto.CreateLinkResponse, error) {
	slug := ""
	if body.Slug != nil {
		slug = *body.Slug
	}

	if slug == "" {
		return ls.createLinkWithGeneratedSlug(ctx, userId, body.OriginalURL)
	}

	link, err := ls.linkRepository.CreateLink(ctx, userId, body.OriginalURL, slug)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateSlug) {
			return dto.CreateLinkResponse{}, ErrDuplicateSlug
		}
		if errors.Is(err, repository.ErrDuplicateOriginalURL) {
			return dto.CreateLinkResponse{}, ErrDuplicateOriginalURL
		}
		return dto.CreateLinkResponse{}, err
	}

	if err := ls.linkCache.SaveOriginalURL(ctx, userId, link.Slug, link.OriginalURL, linkCacheTTL); err != nil {
		return dto.CreateLinkResponse{}, err
	}

	return dto.CreateLinkResponse{
		ID:          link.ID,
		OriginalURL: link.OriginalURL,
		Slug:        link.Slug,
	}, nil
}

func (ls *LinkService) createLinkWithGeneratedSlug(ctx context.Context, userId int, originalURL string) (dto.CreateLinkResponse, error) {
	for range maxSlugGenerationAttempts {
		slug, err := generateRandomSlug(randomSlugLength)
		if err != nil {
			return dto.CreateLinkResponse{}, err
		}

		link, err := ls.linkRepository.CreateLink(ctx, userId, originalURL, slug)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateSlug) {
				continue
			}
			return dto.CreateLinkResponse{}, err
		}

		if err := ls.linkCache.SaveOriginalURL(ctx, userId, link.Slug, link.OriginalURL, linkCacheTTL); err != nil {
			return dto.CreateLinkResponse{}, err
		}

		return dto.CreateLinkResponse{
			ID:          link.ID,
			OriginalURL: link.OriginalURL,
			Slug:        link.Slug,
		}, nil
	}

	return dto.CreateLinkResponse{}, ErrGenerateSlugFailed
}

func generateRandomSlug(length int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var builder strings.Builder
	builder.Grow(length)

	max := big.NewInt(int64(len(alphabet)))
	for range length {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		builder.WriteByte(alphabet[index.Int64()])
	}

	return builder.String(), nil
}
