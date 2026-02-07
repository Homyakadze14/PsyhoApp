package entity

type LoginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"access_token"`
}
