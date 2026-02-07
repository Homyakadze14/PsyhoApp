package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/pkg/postgres"
)

type UserRepository struct {
	postgres.DBConnector
}

func NewUserRepository(pg postgres.DBConnector) *UserRepository {
	return &UserRepository{pg}
}

// Create creates a new user with the default 'user' role
func (r *UserRepository) Create(ctx context.Context, username, password string) (*entity.User, error) {
	const op = "repositories.UserRepository.Create"

	// First get the role ID for 'user' role
	defaultRole := "user"
	var roleID int
	query := `SELECT id FROM role WHERE title = $1 LIMIT 1`
	err := r.QueryRow(ctx, query, defaultRole).Scan(&roleID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get default role: %w", op, err)
	}

	query = `
		INSERT INTO "account"(username, password, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, password, role_id, created_at, updated_at
	`

	var userID int
	var createdAt, updatedAt time.Time

	err = r.QueryRow(ctx, query, username, password, roleID, time.Now(), time.Now()).Scan(
		&userID, &username, &password, &roleID, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user := &entity.User{
		ID:        userID,
		Username:  username,
		Password:  password,
		Role:      defaultRole,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	const op = "repositories.UserRepository.GetByID"

	query := `
		SELECT u.id, u.username, u.password, r.title, u.created_at, u.updated_at
		FROM "account" u
		JOIN role r ON u.role_id = r.id
		WHERE u.id = $1
	`

	var user entity.User
	err := r.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	const op = "repositories.UserRepository.GetByUsername"

	query := `
		SELECT u.id, u.username, u.password, r.title, u.created_at, u.updated_at
		FROM "account" u
		JOIN role r ON u.role_id = r.id
		WHERE u.username = $1
	`

	var user entity.User
	err := r.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	const op = "repositories.UserRepository.Update"

	query := `
		UPDATE "account"
		SET username = $1, password = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.Exec(ctx, query, user.Username, user.Password, time.Now(), user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: user with id %d not found", op, user.ID)
	}

	return nil
}

// Delete removes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	const op = "repositories.UserRepository.Delete"

	query := `DELETE FROM "account" WHERE id = $1`

	result, err := r.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: user with id %d not found", op, id)
	}

	return nil
}

// UpdateUserRole updates a user's role by ID
func (r *UserRepository) UpdateUserRole(ctx context.Context, userID, roleID int) error {
	const op = "repositories.UserRepository.UpdateUserRole"

	query := `UPDATE "account" SET role_id = $1, updated_at = $2 WHERE id = $3`

	result, err := r.Exec(ctx, query, roleID, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: user with id %d not found", op, userID)
	}

	return nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	const op = "repositories.UserRepository.GetAll"

	query := `
		SELECT u.id, u.username, u.password, r.title, u.created_at, u.updated_at
		FROM "account" u
		JOIN role r ON u.role_id = r.id
		ORDER BY u.id
	`

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]entity.User, 0)
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}
