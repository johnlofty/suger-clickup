package models

type User struct {
	ID        int32  `json:"id,omitempty"`
	Password  string `json:"password,omitempty" db:"password"`
	Email     string `json:"email,omitempty" db:"email"`
	CreatedAt int64  `json:"created_at,omitempty" db:"created_at"`
	OrgId     int32  `json:"org_id,omitempty"`
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
