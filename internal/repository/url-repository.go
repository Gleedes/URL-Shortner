package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Gleedes/URL-Shortner/internal/model"
)

type URLRepository interface {
	Save(url *model.URL) error
	GetByShortCode(shortCode string) (*model.URL, error)
	GetByOriginalURL(originalURL string) (*model.URL, error)
}

type sqliteURLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &sqliteURLRepository{db: db}
}

func (sURLRep *sqliteURLRepository) Save(url *model.URL) error {
	_, err := sURLRep.db.Exec("INSERT INTO urls (short_code, original_url) VALUES (?, ?)", url.ShortCode, url.OriginalURL)
	if err != nil {
		return fmt.Errorf("save url: %w", err)
	}
	return nil
}

func (sURLRep *sqliteURLRepository) GetByShortCode(shortCode string) (*model.URL, error) {
	url := &model.URL{}
	err := sURLRep.db.QueryRow("SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = ?", shortCode).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt) // &url. - это адрес поля ... в памяти, просто url - значение поля ...
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrURLNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get by short code: %w", err)
	}

	return url, nil
}

func (sURLRep *sqliteURLRepository) GetByOriginalURL(originalURL string) (*model.URL, error) {
	url := &model.URL{}
	err := sURLRep.db.QueryRow("SELECT id, short_code, original_url, created_at FROM urls WHERE original_url = ?", originalURL).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrURLNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get by original url:  %w", err)
	}
	return url, nil
}
