package dto

import "time"

type TokenResponse struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
