package models

import "time"

type User struct {
	ID        int32     `json:"id,omitempty" db:"user_id"`
	Password  string    `json:"password,omitempty" db:"password"`
	Email     string    `json:"email,omitempty" db:"email"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	OrgId     int32     `json:"org_id,omitempty" db:"org_id"`
}

type RegisterRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
