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
	defaultPage               = 1
	defaultLimit              = 10
)

var (
	ErrDuplicateSlug        = errors.New("slug already used")
	ErrGenerateSlugFailed   = errors.New("failed to generate unique slug")
	ErrDuplicateOriginalURL = errors.New("original URL already used")
	ErrLinkNotFound         = errors.New("link not found")
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

	if err := ls.linkCache.SaveOriginalURL(ctx, link.Slug, link.OriginalURL, linkCacheTTL); err != nil {
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

		if err := ls.linkCache.SaveOriginalURL(ctx, link.Slug, link.OriginalURL, linkCacheTTL); err != nil {
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

func (ls *LinkService) ListLinks(ctx context.Context, userId int, query dto.ListLinksQuery) (dto.ListLinksResponse, error) {
	page := defaultPage
	if query.Page != nil {
		page = *query.Page
	}

	limit := defaultLimit
	if query.Limit != nil {
		limit = *query.Limit
	}

	if page == 0 {
		page = defaultPage
	}
	if limit == 0 {
		limit = defaultLimit
	}

	total, err := ls.linkRepository.CountLinksByUser(ctx, userId)
	if err != nil {
		return dto.ListLinksResponse{}, err
	}

	offset := (page - 1) * limit
	links, err := ls.linkRepository.ListLinksByUser(ctx, userId, limit, offset)
	if err != nil {
		return dto.ListLinksResponse{}, err
	}

	data := make([]dto.CreateLinkResponse, 0, len(links))
	for _, link := range links {
		data = append(data, dto.CreateLinkResponse{
			ID:          link.ID,
			OriginalURL: link.OriginalURL,
			Slug:        link.Slug,
		})
	}

	return dto.ListLinksResponse{
		Data: data,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: calculateTotalPages(total, limit),
		},
	}, nil
}

func (ls *LinkService) DeleteLink(ctx context.Context, userId int, linkId int64) error {
	slug, err := ls.linkRepository.SoftDeleteLink(ctx, userId, linkId)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			return ErrLinkNotFound
		}
		return err
	}

	return ls.linkCache.DeleteOriginalURL(ctx, slug)
}

func (ls *LinkService) ResolveLink(ctx context.Context, slug string) (string, error) {
	originalURL, err := ls.linkCache.GetOriginalURL(ctx, slug)
	if err == nil {
		return originalURL, nil
	}

	if err != nil && !errors.Is(err, repository.ErrLinkCacheMiss) {
		return "", err
	}

	link, err := ls.linkRepository.GetActiveLinkBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			return "", ErrLinkNotFound
		}
		return "", err
	}

	if err := ls.linkCache.SaveOriginalURL(ctx, link.Slug, link.OriginalURL, linkCacheTTL); err != nil {
		return "", err
	}

	return link.OriginalURL, nil
}

func (ls *LinkService) RecordLinkClick(ctx context.Context, slug, ipAddress, userAgent string) error {
	err := ls.linkRepository.InsertLinkClickBySlug(ctx, slug, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, repository.ErrLinkNotFound) {
			return ErrLinkNotFound
		}
		return err
	}

	return nil
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

func calculateTotalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}

	return int((total + int64(limit) - 1) / int64(limit))
}
