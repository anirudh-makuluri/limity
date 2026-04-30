package api

type CheckRequest struct {
	Key    string `json:"key"`
	Limit  int    `json:"limit"`
	Window int    `json:"window"`
}

type CheckResponse struct {
	Allowed   bool  `json:"allowed"`
	Remaining int   `json:"remaining"`
	Reset     int64 `json:"reset"`
}

type TokenClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

type UserProfile struct {
	ID             string `json:"id"`
	ExternalUserID string `json:"external_user_id"`
	Email          string `json:"email"`
	APIKey         string `json:"api_key"`
	CreatedAt      string `json:"created_at"`
}
