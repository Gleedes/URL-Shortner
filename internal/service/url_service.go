package service

import (
	"errors"
	"net/url"
	"time"

	"github.com/Gleedes/URL-Shortner/internal/model"
	"github.com/Gleedes/URL-Shortner/internal/repository"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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
	if err != nil && !errors.Is(err, repository.ErrURLNotFound) {
		return nil, err
	}

	if err == nil {
		return existing, nil
	}

	if errors.Is(err, repository.ErrURLNotFound) {
		url.OriginalURL = originalURL
		url.ShortCode = ""
		url.CreatedAt = time.Now().Format(time.RFC3339)
		if err := urlS.repo.Save(&url); err != nil {
			return nil, err
		}
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
		return false, ErrInvalidScheme
	}

	if u.Host == "" {
		return false, ErrInvalidHost
	}

	return true, nil
}
