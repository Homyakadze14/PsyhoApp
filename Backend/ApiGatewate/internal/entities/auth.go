package entities

import authv1 "github.com/Homyakadze14/PsyhoApp/ApiGatewate/proto/gen/auth"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func (r *RegisterRequest) ToGRPC() *authv1.RegisterRequest {
	return &authv1.RegisterRequest{
		Username: r.Username,
		Password: r.Password,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func (r *LoginRequest) ToGRPC() *authv1.LoginRequest {
	return &authv1.LoginRequest{
		Username: r.Username,
		Password: r.Password,
	}
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

func (r *LogoutRequest) ToGRPC() *authv1.LogoutRequest {
	return &authv1.LogoutRequest{
		AccessToken: r.AccessToken,
	}
}
