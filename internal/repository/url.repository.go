package repository

import (
	"database/sql"
	"time"
	"url-shortener/internal/model"

	"github.com/google/uuid"
)

func SaveURL(db *sql.DB, url model.URL) error {
	query := `
		INSERT INTO urls (id, original, short_code, custom_code, domain, visibility, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.Exec(query,
		url.ID,
		url.Original,
		url.ShortCode,
		url.CustomCode,
		url.Domain,
		url.Visibility,
		url.CreatedAt,
	)
	return err
}

func NewURL(original string, shortCode string) model.URL {
	return model.URL{
		ID:         uuid.New().String(),
		Original:   original,
		ShortCode:  shortCode,
		Visibility: "public",
		CreatedAt:  time.Now(),
	}
}

func GetURLByCode(db *sql.DB, code string) (model.URL, error) {
	query := `SELECT id, original, short_code, custom_code, domain, visibility, created_at FROM urls WHERE short_code = $1`
	row := db.QueryRow(query, code)

	var url model.URL
	err := row.Scan(&url.ID, &url.Original, &url.ShortCode, &url.CustomCode, &url.Domain, &url.Visibility, &url.CreatedAt)
	return url, err

}

func GetAllLinks(db *sql.DB) ([]model.URL, error) {
	rows, err := db.Query("SELECT original, short_code, visibility, created_at FROM urls ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []model.URL
	for rows.Next() {
		var u model.URL
		err := rows.Scan(&u.Original, &u.ShortCode, &u.Visibility, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func UpdateURL(db *sql.DB, url model.URL) error {
	query := `
		UPDATE urls
		SET original = $1, visibility = $2
		WHERE short_code = $3
	`
	_, err := db.Exec(query, url.Original, url.Visibility, url.ShortCode)
	return err
}

func DeleteURLByCode(db *sql.DB, code string) error {
	_, err := db.Exec(`DELETE FROM urls WHERE short_code = $1`, code)
	return err
}
