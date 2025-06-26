package model

import "time"

type Visit struct {
	ID        int       `json:"id"`
	URLID     string    `json:"urlId"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip"`
	Country   string    `json:"country"`
	Browser   string    `json:"browser"`
	Device    string    `json:"device"`
}
