package service

import (
	"database/sql"
	"errors"
	"net/url"
	"time"

	"github.com/Gleedes/URL-Shortner/internal/model"
	"github.com/Gleedes/URL-Shortner/internal/repository"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var ErrInvalidURL = errors.New("invalid URL")

type URLService interface {
	Shorten(originalURL string) (*model.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return urlService{repo: repo}
}

func (urlS urlService) Shorten(originalURL string) (*model.URL, error) {
	if _, err := isValidURL(originalURL); err != nil {
		return nil, err
	}
	url := model.URL{}

	existing, err := urlS.repo.GetByOriginalURL(originalURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err == nil {
		return existing, nil
	}

	if err == sql.ErrNoRows {
		url.OriginalURL = originalURL
		url.ShortCode = ""
		url.CreatedAt = time.Now().Format(time.RFC3339)
		urlS.repo.Save(&url)
		return &url, nil

	}
	return nil, errors.New("unreachable")
}

func isValidURL(originalURL string) (bool, error) {
	u, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return false, ErrInvalidURL
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false, ErrInvalidURL
	}

	if u.Host == "" {
		return false, ErrInvalidURL
	}

	return true, nil
}
