package http

import "time"

type LoginResponse struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Status    string    `json:"status"`
}
