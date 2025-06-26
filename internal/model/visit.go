package model

import "time"

type Visit struct {
	ID        int       `json:"id"`
	URLID     string    `json:"urlId"`
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ip"`
	Country   string    `json:"country"`
	City      string `json:"city"`
	Browser   string    `json:"browser"`
	OS      string `json:"os"`
	Device    string    `json:"device"`
}

//URLID