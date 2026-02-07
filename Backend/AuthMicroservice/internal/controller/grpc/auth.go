package controller

import (
	"context"
	"errors"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	services "github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/usecase"
	authv1 "github.com/Homyakadze14/PsyhoApp/AuthMicroservice/proto/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServiceServer
	auth Auth
}

type Auth interface {
	Login(ctx context.Context, username, password string) (*entity.LoginResponse, error)
	Register(ctx context.Context, username, password string) error
	Logout(ctx context.Context, accessToken string) error
	GenerateAuthCode(ctx context.Context, userID int) (string, error)
	Verify(ctx context.Context, userId int, code string) (bool, error)
	GenerateServiceToken(ctx context.Context, serviceName string) (string, error)
	GetRole(ctx context.Context, userID int) (string, error)
	SetRole(ctx context.Context, userID int, role string) error
	CheckAccessToken(ctx context.Context, accessToken string) (int, error)
	CheckServiceToken(ctx context.Context, serviceToken string) (bool, error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}

// Login implements the login functionality
func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	resp, err := s.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			return nil, status.Error(codes.NotFound, "account not found")
		case errors.Is(err, services.ErrBadCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.LoginResponse{
		Id:          int64(resp.ID),
		AccessToken: resp.Token,
	}, nil
}

// Register implements the registration functionality
func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	err := s.auth.Register(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "account already exists")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.RegisterResponse{
		Success: true,
	}, nil
}

// Logout implements the logout functionality
func (s *serverAPI) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	err := s.auth.Logout(ctx, req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTokenNotFound):
			return nil, status.Error(codes.NotFound, "token not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.LogoutResponse{
		Success: true,
	}, nil
}

// GenerateAuthCode implements the auth code generation functionality
func (s *serverAPI) GenerateAuthCode(ctx context.Context, req *authv1.GenerateAuthCodeRequest) (*authv1.GenerateAuthCodeResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	code, err := s.auth.GenerateAuthCode(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authv1.GenerateAuthCodeResponse{
		Code: code,
	}, nil
}

// Verify implements the verification functionality
func (s *serverAPI) Verify(ctx context.Context, req *authv1.VerifyRequest) (*authv1.VerifyResponse, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	verified, err := s.auth.Verify(ctx, int(req.UserId), req.Code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrVerificationFailed):
			return nil, status.Error(codes.InvalidArgument, "verification failed")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.VerifyResponse{
		Verified: verified,
	}, nil
}

// GenerateServiceToken implements the service token generation functionality
func (s *serverAPI) GenerateServiceToken(ctx context.Context, req *authv1.GenerateServiceTokenRequest) (*authv1.GenerateServiceTokenResponse, error) {
	if req.ServiceName == "" {
		return nil, status.Error(codes.InvalidArgument, "service name is required")
	}

	token, err := s.auth.GenerateServiceToken(ctx, req.ServiceName)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authv1.GenerateServiceTokenResponse{
		Token: token,
	}, nil
}

// GetRole implements the get role functionality
func (s *serverAPI) GetRole(ctx context.Context, req *authv1.GetRoleRequest) (*authv1.GetRoleResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	role, err := s.auth.GetRole(ctx, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authv1.GetRoleResponse{
		Role: role,
	}, nil
}

// SetRole implements the set role functionality
func (s *serverAPI) SetRole(ctx context.Context, req *authv1.SetRoleRequest) (*authv1.SetRoleResponse, error) {
	if req.UserId == 0 || req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and role are required")
	}

	err := s.auth.SetRole(ctx, int(req.UserId), req.Role)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidRole):
			return nil, status.Error(codes.InvalidArgument, "invalid role")
		case errors.Is(err, services.ErrAccountNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.SetRoleResponse{
		Success: true,
	}, nil
}

// CheckAccessToken validates an access token and returns the associated user ID
func (s *serverAPI) CheckAccessToken(ctx context.Context, req *authv1.CheckAccessTokenRequest) (*authv1.CheckAccessTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	userID, err := s.auth.CheckAccessToken(ctx, req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTokenNotFound):
			return nil, status.Error(codes.NotFound, "access token not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &authv1.CheckAccessTokenResponse{
		UserId: int64(userID),
	}, nil
}

// CheckServiceToken validates a service token
func (s *serverAPI) CheckServiceToken(ctx context.Context, req *authv1.CheckServiceTokenRequest) (*authv1.CheckServiceTokenResponse, error) {
	if req.ServiceToken == "" {
		return nil, status.Error(codes.InvalidArgument, "service token is required")
	}

	valid, err := s.auth.CheckServiceToken(ctx, req.ServiceToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &authv1.CheckServiceTokenResponse{
		Valid: valid,
	}, nil
}
