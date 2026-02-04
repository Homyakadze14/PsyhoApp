package entities

import authv1 "github.com/Homyakadze14/PsyhoApp/ApiGatewate/proto/gen/auth"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func (r *RegisterRequest) ToGRPC() *authv1.RegisterRequest {
	return &authv1.RegisterRequest{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func (r *LoginRequest) ToGRPC() *authv1.LoginRequest {
	return &authv1.LoginRequest{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *LogoutRequest) ToGRPC() *authv1.LogoutRequest {
	return &authv1.LogoutRequest{
		RefreshToken: r.RefreshToken,
	}
}

type ActivateAccountRequest struct {
	Link string `json:"link" binding:"required"`
}

func (r *ActivateAccountRequest) ToGRPC() *authv1.ActivateAccountRequest {
	return &authv1.ActivateAccountRequest{
		Link: r.Link,
	}
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (r *RefreshRequest) ToGRPC() *authv1.RefreshRequest {
	return &authv1.RefreshRequest{
		RefreshToken: r.RefreshToken,
	}
}

type SendPasswordLinkRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (r *SendPasswordLinkRequest) ToGRPC() *authv1.SendPasswordLinkRequest {
	return &authv1.SendPasswordLinkRequest{
		Email: r.Email,
	}
}

type ChangePasswordRequest struct {
	Link     string `json:"link" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

func (r *ChangePasswordRequest) ToGRPC() *authv1.ChangePasswordRequest {
	return &authv1.ChangePasswordRequest{
		Password: r.Password,
		Link:     r.Link,
	}
}
