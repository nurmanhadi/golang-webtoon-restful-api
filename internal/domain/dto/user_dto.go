package dto

type UserResponse struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	AvatarFilename string `json:"avatar_filename"`
	AvatarUrl      string `json:"avatar_url"`
}
type UserUpdateUsernameRequest struct {
	Username string `validate:"required,max=100" json:"username"`
}
