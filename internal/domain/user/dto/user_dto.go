package dto

type UserResponse struct {
	Id             string `json:"id"`
	Username       string `json:"username"`
	AvatarFilename string `json:"avatar_file"`
	AvatarUrl      string `json:"avatar_url"`
}
