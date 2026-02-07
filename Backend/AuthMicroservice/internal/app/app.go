package app

import (
	"fmt"
	"log/slog"
	"os"

	grpcapp "github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/app/grpc"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/config"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/entity"
	repository "github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/infra/postgres"
	redisrepo "github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/infra/redis"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/internal/usecase"
	"github.com/Homyakadze14/PsyhoApp/AuthMicroservice/pkg/postgres"
	rds "github.com/Homyakadze14/PsyhoApp/AuthMicroservice/pkg/redis"
)

type App struct {
	db         *postgres.Postgres
	GRPCServer *grpcapp.App
}

func Run(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Database
	pg, err := postgres.New(cfg.Database.URL, postgres.MaxPoolSize(cfg.Database.PoolMax))
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - postgres.New: %w", err).Error())
		os.Exit(1)
	}

	// Redis
	redis, err := rds.New(cfg.Redis)
	if err != nil {
		slog.Error(fmt.Errorf("app - Run - redis.New: %w", err).Error())
		os.Exit(1)
	}

	// Repository
	dbConnector := postgres.NewDBConnector(pg.Pool)
	userRepo := repository.NewUserRepository(dbConnector)
	roleRepo := repository.NewRoleRepository(dbConnector)
	tokenRepo := repository.NewTokenRepository(dbConnector)
	tgConnRepo := repository.NewTgConnectionRepository(dbConnector)
	redisRepo := redisrepo.NewRedisRepository(redis)

	// Usecase
	auth := usecase.NewAuthService(log, userRepo, roleRepo, tokenRepo, tgConnRepo, redisRepo, entity.AuthCode(cfg.AuthCode))

	// GRPC
	gRPCServer := grpcapp.New(log, auth, cfg.GRPC.Port)

	return &App{
		db:         pg,
		GRPCServer: gRPCServer,
	}
}

func (s *App) Shutdown() {
	defer s.db.Close()
	defer s.GRPCServer.Stop()
}
