package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBConnector defines the interface for database operations
type DBConnector interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, args ...interface{}) (int64, error)
}

// Wrapper for pgxpool.Pool to implement DBConnector
type PgxPoolWrapper struct {
	*pgxpool.Pool
}

func NewDBConnector(pool *pgxpool.Pool) DBConnector {
	return &PgxPoolWrapper{pool}
}

func (w *PgxPoolWrapper) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return w.Pool.Query(ctx, query, args...)
}

func (w *PgxPoolWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return w.Pool.QueryRow(ctx, query, args...)
}

func (w *PgxPoolWrapper) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	commandTag, err := w.Pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return commandTag.RowsAffected(), nil
}
