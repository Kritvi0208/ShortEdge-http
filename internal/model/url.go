package model

import "time"

type URL struct {
	ID         string    `json:"id"`
	Original   string    `json:"long_url"`
	ShortCode  string    `json:"code"`
	CustomCode string    `json:"custom_code,omitempty"`
	Domain     string    `json:"domain,omitempty"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"created_at"`
}


func Now() time.Time {
	return time.Now().UTC()
}