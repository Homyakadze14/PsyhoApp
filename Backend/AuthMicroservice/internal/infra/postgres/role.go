package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/pkg/postgres"
)

type RoleRepository struct {
	postgres.DBConnector
}

func NewRoleRepository(pg postgres.DBConnector) *RoleRepository {
	return &RoleRepository{pg}
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, title string) (*entity.Role, error) {
	const op = "repositories.RoleRepository.Create"

	query := `
		INSERT INTO role(title, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, title, created_at, updated_at
	`

	var role entity.Role
	err := r.QueryRow(ctx, query, title, time.Now(), time.Now()).Scan(
		&role.ID, &role.Title, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(ctx context.Context, id int) (*entity.Role, error) {
	const op = "repositories.RoleRepository.GetByID"

	query := `SELECT id, title, created_at, updated_at FROM role WHERE id = $1`

	var role entity.Role
	err := r.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.Title, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

// GetByTitle retrieves a role by title
func (r *RoleRepository) GetByTitle(ctx context.Context, title string) (*entity.Role, error) {
	const op = "repositories.RoleRepository.GetByTitle"

	query := `SELECT id, title, created_at, updated_at FROM role WHERE title = $1`

	var role entity.Role
	err := r.QueryRow(ctx, query, title).Scan(
		&role.ID, &role.Title, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

// Update updates an existing role
func (r *RoleRepository) Update(ctx context.Context, role *entity.Role) error {
	const op = "repositories.RoleRepository.Update"

	query := `
		UPDATE role
		SET title = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.Exec(ctx, query, role.Title, time.Now(), role.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: role with id %d not found", op, role.ID)
	}

	return nil
}

// Delete removes a role by ID
func (r *RoleRepository) Delete(ctx context.Context, id int) error {
	const op = "repositories.RoleRepository.Delete"

	query := `DELETE FROM role WHERE id = $1`

	result, err := r.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result == 0 {
		return fmt.Errorf("%s: role with id %d not found", op, id)
	}

	return nil
}

// GetAll retrieves all roles
func (r *RoleRepository) GetAll(ctx context.Context) ([]entity.Role, error) {
	const op = "repositories.RoleRepository.GetAll"

	query := `SELECT id, title, created_at, updated_at FROM role ORDER BY id`

	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	roles := make([]entity.Role, 0)
	for rows.Next() {
		var role entity.Role
		err := rows.Scan(
			&role.ID, &role.Title, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}
