package auth

type AuthRequest struct {
	Username string `validate:"required,max=100" json:"username"`
	Password string `validate:"required,max=100" json:"password"`
}
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}
