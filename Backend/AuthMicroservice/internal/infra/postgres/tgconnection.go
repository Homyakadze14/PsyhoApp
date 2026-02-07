package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/usecase"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/pkg/postgres"
)

type TgConnectionRepository struct {
	postgres.DBConnector
}

func NewTgConnectionRepository(pg postgres.DBConnector) *TgConnectionRepository {
	return &TgConnectionRepository{pg}
}

// Create creates a new Telegram connection
func (r *TgConnectionRepository) Create(ctx context.Context, userID int, tgUserID int) (*entity.TgConnection, error) {
	const op = "repositories.TgConnectionRepository.Create"

	query := `
		INSERT INTO telegram_connection(user_id, tg_user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, tg_user_id, created_at, updated_at
	`

	var tgConnection entity.TgConnection
	err := r.QueryRow(ctx, query, userID, tgUserID, time.Now(), time.Now()).Scan(
		&tgConnection.ID, &tgConnection.UserID, &tgConnection.TgUserID,
		&tgConnection.CreatedAt, &tgConnection.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &tgConnection, nil
}

// GetByID retrieves a Telegram connection by ID
func (r *TgConnectionRepository) GetByID(ctx context.Context, id int) (*entity.TgConnection, error) {
	const op = "repositories.TgConnectionRepository.GetByID"

	query := `SELECT id, user_id, tg_user_id, created_at, updated_at FROM telegram_connection WHERE id = $1`

	var tgConnection entity.TgConnection
	err := r.QueryRow(ctx, query, id).Scan(
		&tgConnection.ID, &tgConnection.UserID, &tgConnection.TgUserID,
		&tgConnection.CreatedAt, &tgConnection.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &tgConnection, nil
}

// GetByUserID retrieves a Telegram connection by user ID
func (r *TgConnectionRepository) GetByUserID(ctx context.Context, userID int) (*entity.TgConnection, error) {
	const op = "repositories.TgConnectionRepository.GetByUserID"

	query := `SELECT id, user_id, tg_user_id, created_at, updated_at FROM telegram_connection WHERE user_id = $1`

	var tgConnection entity.TgConnection
	err := r.QueryRow(ctx, query, userID).Scan(
		&tgConnection.ID, &tgConnection.UserID, &tgConnection.TgUserID,
		&tgConnection.CreatedAt, &tgConnection.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrTgConnNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &tgConnection, nil
}

// GetByTgUserID retrieves a Telegram connection by Telegram user ID
func (r *TgConnectionRepository) GetByTgUserID(ctx context.Context, tgUserID int) (*entity.TgConnection, error) {
	const op = "repositories.TgConnectionRepository.GetByTgUserID"

	query := `SELECT id, user_id, tg_user_id, created_at, updated_at FROM telegram_connection WHERE tg_user_id = $1`

	var tgConnection entity.TgConnection
	err := r.QueryRow(ctx, query, tgUserID).Scan(
		&tgConnection.ID, &tgConnection.UserID, &tgConnection.TgUserID,
		&tgConnection.CreatedAt, &tgConnection.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.ErrTgConnNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &tgConnection, nil
}

// Update updates an existing Telegram connection
func (r *TgConnectionRepository) Update(ctx context.Context, tgConnection *entity.TgConnection) error {
	const op = "repositories.TgConnectionRepository.Update"

	query := `
		UPDATE telegram_connection
		SET user_id = $1, tg_user_id = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.Exec(ctx, query, tgConnection.UserID, tgConnection.TgUserID, time.Now(), tgConnection.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: telegram connection with id %d not found", op, tgConnection.ID)
	}

	return nil
}

// Delete removes a Telegram connection by ID
func (r *TgConnectionRepository) Delete(ctx context.Context, id int) error {
	const op = "repositories.TgConnectionRepository.Delete"

	query := `DELETE FROM telegram_connection WHERE id = $1`

	result, err := r.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: telegram connection with id %d not found", op, id)
	}

	return nil
}

// GetAll retrieves all Telegram connections
func (r *TgConnectionRepository) GetAll(ctx context.Context) ([]entity.TgConnection, error) {
	const op = "repositories.TgConnectionRepository.GetAll"

	query := `SELECT id, user_id, tg_user_id, created_at, updated_at FROM telegram_connection ORDER BY id`

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	tgConnections := make([]entity.TgConnection, 0)
	for rows.Next() {
		var tgConnection entity.TgConnection
		err := rows.Scan(
			&tgConnection.ID, &tgConnection.UserID, &tgConnection.TgUserID,
			&tgConnection.CreatedAt, &tgConnection.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tgConnections = append(tgConnections, tgConnection)
	}

	return tgConnections, nil
}

// DeleteByUserID removes a Telegram connection by user ID
func (r *TgConnectionRepository) DeleteByUserID(ctx context.Context, userID int) error {
	const op = "repositories.TgConnectionRepository.DeleteByUserID"

	query := `DELETE FROM telegram_connection WHERE user_id = $1`

	result, err := r.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: telegram connection with user id %d not found", op, userID)
	}

	return nil
}

// DeleteByTgUserID removes a Telegram connection by Telegram user ID
func (r *TgConnectionRepository) DeleteByTgUserID(ctx context.Context, tgUserID int) error {
	const op = "repositories.TgConnectionRepository.DeleteByTgUserID"

	query := `DELETE FROM telegram_connection WHERE tg_user_id = $1`

	result, err := r.Exec(ctx, query, tgUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: telegram connection with tg user id %v not found", op, tgUserID)
	}

	return nil
}
