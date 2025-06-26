package repository

import "database/sql"

func ShortCodeExists(db *sql.DB, code string) (bool, error) {
	row := db.QueryRow("SELECT 1 FROM urls WHERE short_code = $1", code)
	var dummy int
	err := row.Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
