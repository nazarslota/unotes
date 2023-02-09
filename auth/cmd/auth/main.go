package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/nazarslota/unotes/auth/internal/config"
	"github.com/nazarslota/unotes/auth/internal/handler/grpc"
	"github.com/nazarslota/unotes/auth/internal/handler/rest"
	"github.com/nazarslota/unotes/auth/internal/service"
	"github.com/nazarslota/unotes/auth/internal/storage"
	"github.com/nazarslota/unotes/auth/internal/storage/postgres"
	"github.com/nazarslota/unotes/auth/internal/storage/redis"
	"github.com/nazarslota/unotes/auth/pkg/jwt"
	"github.com/nazarslota/unotes/auth/pkg/logger"
	"github.com/nazarslota/unotes/auth/pkg/utils"
)

var log logger.Logger

func init() {
	var logs io.Writer
	logs, err := os.OpenFile(config.C().Auth.Log, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logs = io.Discard
	}
	out := io.MultiWriter(logs, logger.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen})
	log = logger.NewLogger(out).With().Timestamp().Logger()
}

func main() {
	log.Info("Connecting to the PostgreSQL database...")
	postgresDB, err := postgres.NewPostgreSQL(context.Background(), postgres.Config{
		Host:     config.C().PostgreSQL.Host,
		Port:     config.C().PostgreSQL.Port,
		Username: config.C().PostgreSQL.Username,
		Password: config.C().PostgreSQL.Password,
		DBName:   config.C().PostgreSQL.DBName,
		SSLMode:  config.C().PostgreSQL.SSLMode,
	})
	if err != nil {
		log.FatalFields("Failed to connect to PostgreSQL database.", map[string]any{"error": err})
	} else {
		log.Info("Successfully connected to PostgreSQL database.")
	}

	log.Info("Connecting to the Redis database...")
	redisDB, err := redis.NewRedis(context.Background(), redis.Config{
		Addr:     config.C().Redis.Addr,
		Password: config.C().Redis.Password,
		DB:       config.C().Redis.DB,
	})
	if err != nil {
		log.FatalFields("Failed to connect to Redis database.", map[string]any{"error": err})
	} else {
		log.Info("Successfully connected to Redis database.")
	}

	repositories := storage.NewRepositoryProvider(
		storage.WithPostgreSQLUserRepository(postgresDB),
		storage.WithRedisRefreshTokenRepository(redisDB),
	)

	services := service.NewServices(service.OAuth2ServiceOptions{
		AccessTokenManager:     jwt.NewAccessTokenManagerHMAC(config.C().Auth.AccessTokenSecret),
		AccessTokenExpiresIn:   config.C().Auth.AccessTokenExpiresIn,
		RefreshTokenManager:    jwt.NewRefreshTokenManagerHMAC(config.C().Auth.RefreshTokenSecret),
		RefreshTokenExpiresIn:  config.C().Auth.RefreshTokenExpiresIn,
		UserRepository:         repositories.UserRepository,
		RefreshTokenRepository: repositories.RefreshTokenRepository,
	})

	restAddress := net.JoinHostPort(config.C().Auth.HostREST, config.C().Auth.PortREST)
	restServer := rest.NewHandler(
		rest.WithServices(services),
		rest.WithAddress(restAddress),
		rest.WithLogger(log),
		rest.WithDebug(config.C().Auth.Debug),
	).Server()

	log.InfoFields("Starting a REST server...", map[string]any{"address": restAddress})
	go func() {
		if err := restServer.Serve(); err != nil {
			log.FatalFields("Error occurred while running REST server.", map[string]any{"error": err})
		}
	}()

	time.Sleep(time.Second)
	log.Info("The REST server is successfully started.")

	grpcAddress := net.JoinHostPort(config.C().Auth.HostGRPC, config.C().Auth.PortGRPC)
	grpcServer := grpc.NewHandler(
		grpc.WithService(services),
		grpc.WithAddress(grpcAddress),
		grpc.WithLogger(log),
	).Server()

	log.InfoFields("Starting a gRPC server...", map[string]any{"address": grpcAddress})
	go func() {
		if err := grpcServer.Serve(); err != nil {
			log.FatalFields("Error occurred while running gRPC server.", map[string]any{"error": err})
		}
	}()

	time.Sleep(time.Second)
	log.Info("The gRPC server is successfully started.")

	<-utils.GracefulShutdown()
	log.Info("Shutdown of the REST server...")
	if err := restServer.Shutdown(context.Background()); err != nil {
		log.ErrorFields("Error during REST server shutdown.", map[string]any{"error": err})
	} else {
		log.Info("REST server was successfully shut down.")
	}

	log.Info("Shutdown of the gRPC server...")
	if err := grpcServer.Shutdown(context.Background()); err != nil {
		log.ErrorFields("Error during gRPC server shutdown.", map[string]any{"error": err})
	} else {
		log.Info("gRPC server was successfully shut down.")
	}

	log.Info("Closing the connection to PostgreSQL...")
	if err := postgresDB.Close(); err != nil {
		log.ErrorFields("Failed to close connection to PostgreSQL.", map[string]any{"error": err})
	} else {
		log.Info("The connection to PostgreSQL is successfully closed.")
	}

	log.Info("Closing the connection to Redis...")
	if err := redisDB.Close(); err != nil {
		log.ErrorFields("Failed to close connection to Redis.", map[string]any{"error": err})
	} else {
		log.Info("The connection to Redis is successfully closed.")
	}
}
