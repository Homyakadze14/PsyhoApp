package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/big"
	"time"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAccountAlreadyExists = errors.New("account with this credentials already exists")
	ErrAccountNotFound      = errors.New("account with this credentials not found")
	ErrBadCredentials       = errors.New("bad credentials")
	ErrTokenNotFound        = errors.New("token not found")
	ErrLinkNotFound         = errors.New("link not found")
	ErrNotActivated         = errors.New("not activated account")
	ErrInvalidRole          = errors.New("invalid role")
	ErrVerificationFailed   = errors.New("verification failed")
	ErrTgConnNotFound       = errors.New("telegram connectino not found")
	ErrCacheNotFound        = errors.New("cache not found")
)

type UserRepoI interface {
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Create(ctx context.Context, username, password string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
	UpdateUserRole(ctx context.Context, userID, roleID int) error
}

type RoleRepoI interface {
	GetByID(ctx context.Context, id int) (*entity.Role, error)
	GetByTitle(ctx context.Context, title string) (*entity.Role, error)
	Create(ctx context.Context, title string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id int) error
}

type TokenRepoI interface {
	CreateServiceToken(ctx context.Context, serviceName, token string) (*entity.SerivceToken, error)
	GetServiceTokenByServiceName(ctx context.Context, serviceName string) (*entity.SerivceToken, error)
	GetServiceTokenByToken(ctx context.Context, token string) (*entity.SerivceToken, error)
	CreateAccessToken(ctx context.Context, userID int, token string) (*entity.AccessToken, error)
	GetAccessTokenByToken(ctx context.Context, token string) (*entity.AccessToken, error)
	DeleteAccessToken(ctx context.Context, id int) error
}

type TgConnectionRepoI interface {
	Create(ctx context.Context, userID int, tgUserID int) (*entity.TgConnection, error)
	GetByUserID(ctx context.Context, userID int) (*entity.TgConnection, error)
}

type RedisRepository interface {
	Set(ctx context.Context, key string, value any, expTime time.Duration) error
	Del(ctx context.Context, key string) (res int64, err error)
	Get(ctx context.Context, key string, dest any) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

type AuthService struct {
	log       *slog.Logger
	userRepo  UserRepoI
	roleRepo  RoleRepoI
	tokenRepo TokenRepoI
	tgConn    TgConnectionRepoI
	authCodes RedisRepository
	authCode  entity.AuthCode
}

func NewAuthService(
	log *slog.Logger,
	userRepo UserRepoI,
	roleRepo RoleRepoI,
	tokenRepo TokenRepoI,
	tgConn TgConnectionRepoI,
	authCodes RedisRepository,
	authCode entity.AuthCode,
) *AuthService {
	return &AuthService{
		log:       log,
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		tokenRepo: tokenRepo,
		tgConn:    tgConn,
		authCodes: authCodes,
		authCode:  authCode,
	}
}

// Login handles user login
func (s *AuthService) Login(ctx context.Context, username, password string) (*entity.LoginResponse, error) {
	const op = "AuthService.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("login attempt")

	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return nil, ErrAccountNotFound
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Error("invalid password")
		return nil, ErrBadCredentials
	}

	// Generate access token
	accessToken, err := s.generateAccessToken()
	if err != nil {
		log.Error("failed to generate access token", slog.String("error", err.Error()))
		return nil, err
	}

	// Store access token in database
	_, err = s.tokenRepo.CreateAccessToken(ctx, user.ID, accessToken)
	if err != nil {
		log.Error("failed to store access token", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("login successful")
	return &entity.LoginResponse{
		ID:    int(user.ID),
		Token: accessToken,
	}, nil
}

// Register handles user registration
func (s *AuthService) Register(ctx context.Context, username, password string) error {
	const op = "AuthService.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	log.Info("registration attempt")

	// Check if user already exists
	_, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		log.Error("user already exists")
		return ErrAccountAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return err
	}

	// Create user (with default 'user' role)
	_, err = s.userRepo.Create(ctx, username, string(hashedPassword))
	if err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("registration successful")
	return nil
}

// Logout handles user logout
func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	const op = "AuthService.Logout"

	log := s.log.With(
		slog.String("op", op),
		slog.String("token", accessToken),
	)

	log.Info("logout attempt")

	// Find access token
	token, err := s.tokenRepo.GetAccessTokenByToken(ctx, accessToken)
	if err != nil {
		log.Error("token not found", slog.String("error", err.Error()))
		return ErrTokenNotFound
	}

	// Delete the token from database
	err = s.tokenRepo.DeleteAccessToken(ctx, token.ID)
	if err != nil {
		log.Error("failed to delete token", slog.String("error", err.Error()))
		return err
	}

	log.Info("logout successful", slog.Int("user_id", token.UserID))
	return nil
}

// GenerateAuthCode generates an authentication code for a user
func (s *AuthService) GenerateAuthCode(ctx context.Context, userID int) (string, error) {
	const op = "AuthService.GenerateAuthCode"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("user_id", userID),
	)

	log.Info("generating auth code")

	// Generate random code
	code, err := s.generateRandomCode()
	if err != nil {
		log.Error("failed to generate auth code", slog.String("error", err.Error()))
		return "", err
	}

	// Insert to cache
	go func() {
		err := s.authCodes.Set(ctx, code, userID, s.authCode.TTL)
		log.Error("failed to set auth code in cache", slog.String("error", err.Error()))
	}()

	log.Info("auth code generated", slog.String("code", code))
	return code, nil
}

// Verify verifies an authentication code
func (s *AuthService) Verify(ctx context.Context, userID int, code string) (bool, error) {
	const op = "AuthService.Verify"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("user_id", userID),
		slog.String("code", code),
	)

	log.Info("verifying auth code")

	// Check if code exists
	var tgUserID int
	err := s.authCodes.Get(ctx, code, &tgUserID)
	if err != nil {
		switch {
		case errors.Is(err, ErrCacheNotFound):
			log.Error("auth code not found")
			return false, ErrVerificationFailed
		default:
			log.Error("get from cache error", slog.String("error", err.Error()))
			return false, err
		}
	}

	// Compare ids
	tgConn, err := s.tgConn.GetByUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrTgConnNotFound):
			tgConn, err = s.tgConn.Create(ctx, userID, tgUserID)
			if err != nil {
				log.Error("failed to create tg connection", slog.String("error", err.Error()))
				return false, ErrVerificationFailed
			}
		default:
			log.Error("failed to get tg connection", slog.String("error", err.Error()))
			return false, ErrVerificationFailed
		}
	}

	if tgUserID != tgConn.TgUserID {
		log.Error("invalid user_id")
		return false, ErrVerificationFailed
	}

	_, err = s.authCodes.Del(ctx, code)
	if err != nil {
		log.Error("failed to delete code from cache", slog.String("error", err.Error()))
		return false, err
	}

	log.Info("auth code verified")
	return true, nil
}

// GenerateServiceToken generates a service token for a service
func (s *AuthService) GenerateServiceToken(ctx context.Context, serviceName string) (string, error) {
	const op = "AuthService.GenerateServiceToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("service_name", serviceName),
	)

	log.Info("generating service token")

	// Check if service token already exists
	existingToken, err := s.tokenRepo.GetServiceTokenByServiceName(ctx, serviceName)
	if err == nil {
		// If token exists, return it
		log.Info("service token already exists")
		return existingToken.Token, nil
	}

	// Generate new token
	serviceToken, err := s.generateServiceToken()
	if err != nil {
		log.Error("failed to generate service token", slog.String("error", err.Error()))
		return "", err
	}

	// Store service token in database
	_, err = s.tokenRepo.CreateServiceToken(ctx, serviceName, serviceToken)
	if err != nil {
		log.Error("failed to store service token", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("service token generated")
	return serviceToken, nil
}

// CheckAccessToken validates an access token and returns the associated user ID
func (s *AuthService) CheckAccessToken(ctx context.Context, accessToken string) (int, error) {
	const op = "AuthService.CheckAccessToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("token", accessToken),
	)

	log.Info("checking access token")

	// Find access token in database
	token, err := s.tokenRepo.GetAccessTokenByToken(ctx, accessToken)
	if err != nil {
		log.Error("access token not found", slog.String("error", err.Error()))
		return 0, ErrTokenNotFound
	}

	log.Info("access token validated", slog.Int("user_id", token.UserID))
	return token.UserID, nil
}

// CheckServiceToken validates a service token
func (s *AuthService) CheckServiceToken(ctx context.Context, serviceToken string) (bool, error) {
	const op = "AuthService.CheckServiceToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("token", serviceToken),
	)

	log.Info("checking service token")

	// Find service token in database
	_, err := s.tokenRepo.GetServiceTokenByToken(ctx, serviceToken)
	if err != nil {
		log.Error("service token not found", slog.String("error", err.Error()))
		return false, nil
	}

	log.Info("service token validated")
	return true, nil
}

// GetRole gets the role of a user
func (s *AuthService) GetRole(ctx context.Context, userID int) (string, error) {
	const op = "AuthService.GetRole"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("user_id", userID),
	)

	log.Info("getting user role")

	// Get user by ID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return "", ErrAccountNotFound
	}

	log.Info("user role retrieved", slog.String("role", user.Role))
	return user.Role, nil
}

// SetRole sets the role of a user
func (s *AuthService) SetRole(ctx context.Context, userID int, role string) error {
	const op = "AuthService.SetRole"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("user_id", userID),
		slog.String("role", role),
	)

	log.Info("setting user role")

	// Validate role exists
	roleEntity, err := s.roleRepo.GetByTitle(ctx, role)
	if err != nil {
		log.Error("role does not exist", slog.String("error", err.Error()))
		return ErrInvalidRole
	}

	// Get user by ID to verify it exists
	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		return ErrAccountNotFound
	}

	// Update user's role
	err = s.userRepo.UpdateUserRole(ctx, userID, roleEntity.ID)
	if err != nil {
		log.Error("failed to update user role", slog.String("error", err.Error()))
		return err
	}

	log.Info("user role set")
	return nil
}

// generateAccessToken generates a random access token
func (s *AuthService) generateAccessToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateServiceToken generates a random service token
func (s *AuthService) generateServiceToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateRandomCode generates a random authentication code consisting of digits only
func (s *AuthService) generateRandomCode() (string, error) {
	code := make([]byte, s.authCode.Length)

	// Fill the byte slice with random digits (0-9)
	for i := range s.authCode.Length {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code[i] = byte('0' + num.Int64())
	}

	return string(code), nil
}
