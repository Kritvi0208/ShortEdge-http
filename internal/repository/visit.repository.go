package repository

import (
	"database/sql"
	"url-shortener/internal/model"
)

func SaveVisit(db *sql.DB, visit model.Visit) error {
	query := `
		INSERT INTO visits (url_id, timestamp, ip_address, country, browser, device)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(query,
		visit.URLID,
		visit.Timestamp,
		visit.IPAddress,
		visit.Country,
		visit.Browser,
		visit.Device,
	)
	return err
}


func GetVisitsByURLID(db *sql.DB, urlID string) ([]model.Visit, error) {
	query := `
		SELECT id, timestamp, ip_address, country, browser, device
		FROM visits
		WHERE url_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := db.Query(query, urlID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []model.Visit
	for rows.Next() {
		var v model.Visit
		v.URLID = urlID
		err := rows.Scan(&v.ID, &v.Timestamp, &v.IPAddress, &v.Country, &v.Browser, &v.Device)
		if err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	return visits, nil
}

// package repository

// import (
// 	"database/sql"
// 	"url-shortener/internal/model"
// )

// func SaveVisit(db *sql.DB, visit model.Visit) error {
// 	query := `
// 		INSERT INTO visits (url_id, ip_address, country, timestamp, browser, device)
// VALUES ($1, $2, $3, $4, $5, $6)
// 	`
// 	_, err := db.Exec(query,
// 		visit.URLID,     // $1 → url_id
// 		visit.IPAddress, // $2 → ip_address
// 		visit.Country,   // $3 → country
// 		visit.Timestamp, // $4 → timestamp
// 		visit.Browser,   // $5 → browser
// 		visit.Device,    // $6 → device
// 	)

// 	return err
// }

// func GetVisitsByURLID(db *sql.DB, urlID string) ([]model.Visit, error) {
// 	query := `
// 		SELECT id, timestamp, ip_address, country, browser, device
// 		FROM visits
// 		WHERE url_id = $1
// 		ORDER BY timestamp DESC
// 	`

// 	rows, err := db.Query(query, urlID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var visits []model.Visit
// 	for rows.Next() {
// 		var v model.Visit
// 		v.URLID = urlID
// 		err := rows.Scan(&v.ID, &v.Timestamp, &v.IPAddress, &v.Country, &v.Browser, &v.Device)
// 		if err != nil {
// 			return nil, err
// 		}
// 		visits = append(visits, v)
// 	}
// 	return visits, nil
// }
