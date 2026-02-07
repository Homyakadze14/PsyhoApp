package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/pkg/postgres"
)

type TokenRepository struct {
	postgres.DBConnector
}

func NewTokenRepository(pg postgres.DBConnector) *TokenRepository {
	return &TokenRepository{pg}
}

// CreateServiceToken creates a new service token
func (r *TokenRepository) CreateServiceToken(ctx context.Context, serviceName, token string) (*entity.SerivceToken, error) {
	const op = "repositories.TokenRepository.CreateServiceToken"

	query := `
		INSERT INTO service_token(service_name, token, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, service_name, token, created_at, updated_at
	`

	var serviceToken entity.SerivceToken
	err := r.QueryRow(ctx, query, serviceName, token, time.Now(), time.Now()).Scan(
		&serviceToken.ID, &serviceToken.ServiceName, &serviceToken.Token,
		&serviceToken.CreatedAt, &serviceToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &serviceToken, nil
}

// GetServiceTokenByID retrieves a service token by ID
func (r *TokenRepository) GetServiceTokenByID(ctx context.Context, id int) (*entity.SerivceToken, error) {
	const op = "repositories.TokenRepository.GetServiceTokenByID"

	query := `SELECT id, service_name, token, created_at, updated_at FROM service_token WHERE id = $1`

	var serviceToken entity.SerivceToken
	err := r.QueryRow(ctx, query, id).Scan(
		&serviceToken.ID, &serviceToken.ServiceName, &serviceToken.Token,
		&serviceToken.CreatedAt, &serviceToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &serviceToken, nil
}

// GetServiceTokenByServiceName retrieves a service token by service name
func (r *TokenRepository) GetServiceTokenByServiceName(ctx context.Context, serviceName string) (*entity.SerivceToken, error) {
	const op = "repositories.TokenRepository.GetServiceTokenByServiceName"

	query := `SELECT id, service_name, token, created_at, updated_at FROM service_token WHERE service_name = $1`

	var serviceToken entity.SerivceToken
	err := r.QueryRow(ctx, query, serviceName).Scan(
		&serviceToken.ID, &serviceToken.ServiceName, &serviceToken.Token,
		&serviceToken.CreatedAt, &serviceToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &serviceToken, nil
}

// GetServiceTokenByToken retrieves a service token by token value
func (r *TokenRepository) GetServiceTokenByToken(ctx context.Context, token string) (*entity.SerivceToken, error) {
	const op = "repositories.TokenRepository.GetServiceTokenByToken"

	query := `SELECT id, service_name, token, created_at, updated_at FROM service_token WHERE token = $1`

	var serviceToken entity.SerivceToken
	err := r.QueryRow(ctx, query, token).Scan(
		&serviceToken.ID, &serviceToken.ServiceName, &serviceToken.Token,
		&serviceToken.CreatedAt, &serviceToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &serviceToken, nil
}

// UpdateServiceToken updates an existing service token
func (r *TokenRepository) UpdateServiceToken(ctx context.Context, serviceToken *entity.SerivceToken) error {
	const op = "repositories.TokenRepository.UpdateServiceToken"

	query := `
		UPDATE service_token
		SET service_name = $1, token = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.Exec(ctx, query, serviceToken.ServiceName, serviceToken.Token, time.Now(), serviceToken.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: service token with id %d not found", op, serviceToken.ID)
	}

	return nil
}

// DeleteServiceToken removes a service token by ID
func (r *TokenRepository) DeleteServiceToken(ctx context.Context, id int) error {
	const op = "repositories.TokenRepository.DeleteServiceToken"

	query := `DELETE FROM service_token WHERE id = $1`

	result, err := r.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: service token with id %d not found", op, id)
	}

	return nil
}

// GetAllServiceTokens retrieves all service tokens
func (r *TokenRepository) GetAllServiceTokens(ctx context.Context) ([]entity.SerivceToken, error) {
	const op = "repositories.TokenRepository.GetAllServiceTokens"

	query := `SELECT id, service_name, token, created_at, updated_at FROM service_token ORDER BY id`

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	serviceTokens := make([]entity.SerivceToken, 0)
	for rows.Next() {
		var serviceToken entity.SerivceToken
		err := rows.Scan(
			&serviceToken.ID, &serviceToken.ServiceName, &serviceToken.Token,
			&serviceToken.CreatedAt, &serviceToken.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		serviceTokens = append(serviceTokens, serviceToken)
	}

	return serviceTokens, nil
}

// CreateAccessToken creates a new access token
func (r *TokenRepository) CreateAccessToken(ctx context.Context, userID int, token string) (*entity.AccessToken, error) {
	const op = "repositories.TokenRepository.CreateAccessToken"

	query := `
		INSERT INTO token(user_id, access_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, access_token, created_at, updated_at
	`

	var accessToken entity.AccessToken
	err := r.QueryRow(ctx, query, userID, token, time.Now(), time.Now()).Scan(
		&accessToken.ID, &accessToken.UserID, &accessToken.Token,
		&accessToken.CreatedAt, &accessToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &accessToken, nil
}

// GetAccessTokenByID retrieves an access token by ID
func (r *TokenRepository) GetAccessTokenByID(ctx context.Context, id int) (*entity.AccessToken, error) {
	const op = "repositories.TokenRepository.GetAccessTokenByID"

	query := `SELECT id, user_id, access_token, created_at, updated_at FROM token WHERE id = $1`

	var accessToken entity.AccessToken
	err := r.QueryRow(ctx, query, id).Scan(
		&accessToken.ID, &accessToken.UserID, &accessToken.Token,
		&accessToken.CreatedAt, &accessToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &accessToken, nil
}

// GetAccessTokenByToken retrieves an access token by token value
func (r *TokenRepository) GetAccessTokenByToken(ctx context.Context, token string) (*entity.AccessToken, error) {
	const op = "repositories.TokenRepository.GetAccessTokenByToken"

	query := `SELECT id, user_id, access_token, created_at, updated_at FROM token WHERE access_token = $1`

	var accessToken entity.AccessToken
	err := r.QueryRow(ctx, query, token).Scan(
		&accessToken.ID, &accessToken.UserID, &accessToken.Token,
		&accessToken.CreatedAt, &accessToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &accessToken, nil
}

// GetAccessTokenByUserID retrieves an access token by user ID
func (r *TokenRepository) GetAccessTokenByUserID(ctx context.Context, userID int) (*entity.AccessToken, error) {
	const op = "repositories.TokenRepository.GetAccessTokenByUserID"

	query := `SELECT id, user_id, access_token, created_at, updated_at FROM token WHERE user_id = $1`

	var accessToken entity.AccessToken
	err := r.QueryRow(ctx, query, userID).Scan(
		&accessToken.ID, &accessToken.UserID, &accessToken.Token,
		&accessToken.CreatedAt, &accessToken.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &accessToken, nil
}

// UpdateAccessToken updates an existing access token
func (r *TokenRepository) UpdateAccessToken(ctx context.Context, accessToken *entity.AccessToken) error {
	const op = "repositories.TokenRepository.UpdateAccessToken"

	query := `
		UPDATE token
		SET user_id = $1, access_token = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.Exec(ctx, query, accessToken.UserID, accessToken.Token, time.Now(), accessToken.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: access token with id %d not found", op, accessToken.ID)
	}

	return nil
}

// DeleteAccessToken removes an access token by ID
func (r *TokenRepository) DeleteAccessToken(ctx context.Context, id int) error {
	const op = "repositories.TokenRepository.DeleteAccessToken"

	query := `DELETE FROM token WHERE id = $1`

	result, err := r.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: access token with id %d not found", op, id)
	}

	return nil
}

// GetAllAccessTokens retrieves all access tokens
func (r *TokenRepository) GetAllAccessTokens(ctx context.Context) ([]entity.AccessToken, error) {
	const op = "repositories.TokenRepository.GetAllAccessTokens"

	query := `SELECT id, user_id, access_token, created_at, updated_at FROM token ORDER BY id`

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	accessTokens := make([]entity.AccessToken, 0)
	for rows.Next() {
		var accessToken entity.AccessToken
		err := rows.Scan(
			&accessToken.ID, &accessToken.UserID, &accessToken.Token,
			&accessToken.CreatedAt, &accessToken.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		accessTokens = append(accessTokens, accessToken)
	}

	return accessTokens, nil
}
