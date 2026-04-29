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

type APIKey struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Key       string  `json:"key"`
	CreatedAt string  `json:"created_at"`
	RevokedAt *string `json:"revoked_at"`
}

type TokenClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}
