package services

import (
	"fmt"
	"log/slog"

	"github.com/Homyakadze14/PsyhoApp/ApiGatewate/internal/config"
	authv1 "github.com/Homyakadze14/PsyhoApp/ApiGatewate/proto/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	log  *slog.Logger
	cfg  config.AuthServiceConfig
	conn *grpc.ClientConn
}

func NewAuthService(log *slog.Logger, cfg config.AuthServiceConfig) *AuthService {
	return &AuthService{
		log: log,
		cfg: cfg,
	}
}

func (s *AuthService) Connect() authv1.AuthClient {
	const op = "AuthService.Connect"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("trying to connect to auth service")
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(s.cfg.Addr, opts...)
	if err != nil {
		log.Error("failed to connect to auth service")
		panic(fmt.Errorf("%s: %w", op, err))
	}
	s.conn = conn

	client := authv1.NewAuthClient(conn)
	log.Info("successfully connected to the auth service")

	return client
}

func (s *AuthService) CloseConn() error {
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
