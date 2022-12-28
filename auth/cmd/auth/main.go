package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	redisdriver "github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/udholdenhed/unotes/auth/internal/config"
	handlerrest "github.com/udholdenhed/unotes/auth/internal/handler/grpc"
	handlerhttp "github.com/udholdenhed/unotes/auth/internal/handler/rest"
	"github.com/udholdenhed/unotes/auth/internal/service"
	"github.com/udholdenhed/unotes/auth/internal/storage"
	"github.com/udholdenhed/unotes/auth/internal/storage/postgres"
	"github.com/udholdenhed/unotes/auth/internal/storage/redis"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	InitLogger()

	//postgresDB, err := NewPostgreSQL()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Filed to connect to PostgreSQL.")
	//}
	//
	//redisDB, err := NewRedisDB()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Filed to connect to Redis.")
	//}

	repositories := storage.NewRepositoryProvider(
		storage.WithMemoryUserRepository(),
		storage.WithMemoryRefreshTokenRepository(),
	)

	services := service.NewService(&service.OAuth2ServiceOptions{
		AccessTokenSecret:      config.C().Auth.AccessTokenSecret,
		RefreshTokenSecret:     config.C().Auth.RefreshTokenSecret,
		AccessTokenExpiresIn:   config.C().Auth.AccessTokenExpiresIn,
		RefreshTokenExpiresIn:  config.C().Auth.RefreshTokenExpiresIn,
		UserRepository:         repositories.UserRepository,
		RefreshTokenRepository: repositories.RefreshTokenRepository,
	})

	addrHTTP := net.JoinHostPort(
		config.C().Auth.HostHTTP,
		config.C().Auth.PortHTTP,
	)
	serverHTTP := handlerhttp.NewHandler(services, &log.Logger).Init(addrHTTP)

	go func() {
		if err := serverHTTP.ListenAndServe(); err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
			default:
				log.Fatal().Err(err).Msg("Error occurred while running HTTP server.")
			}
		}
	}()
	log.Info().Msg("HTTP server started.")

	addrGRPC := net.JoinHostPort(
		config.C().Auth.HostGRPC,
		config.C().Auth.PortGRPC,
	)
	serverGRPC := handlerrest.NewHandler(services, &log.Logger).Init(addrGRPC)

	go func() {
		listener, err := net.Listen("tcp", addrGRPC)
		if err != nil {
			log.Fatal().Err(err).Msg("Error occurred while running gRPC server.")
		}

		if err := serverGRPC.Serve(listener); err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
			default:
				log.Fatal().Err(err).Msg("Error occurred while running HTTP server.")
			}
		}
	}()
	log.Info().Msg("gRPC server started.")

	<-utils.GracefulShutdown()
	log.Info().Msg("Shutting down.")

	ShutdownHTTPServer(serverHTTP)
	ShutdownGRPCServer(serverGRPC)

	//ClosePostgreSQL(postgresDB)
	//CloseRedisDB(redisDB)

	log.Info().Msg("Shutdown completed.")
}

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC1123,
	})
}

func NewPostgreSQL() (*sqlx.DB, error) {
	p, err := postgres.NewPostgreSQL(context.Background(), &postgres.Config{
		Host:     config.C().PostgreSQL.Host,
		Port:     config.C().PostgreSQL.Port,
		Username: config.C().PostgreSQL.Username,
		Password: config.C().PostgreSQL.Password,
		DBName:   config.C().PostgreSQL.DBName,
		SSLMode:  config.C().PostgreSQL.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("main: %w", err)
	}
	return p, nil
}

func NewRedisDB() (*redisdriver.Client, error) {
	r, err := redis.NewRedis(context.Background(), &redis.Config{
		Addr:     config.C().Redis.Addr,
		Password: config.C().Redis.Password,
		DB:       config.C().Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("main: %w", err)
	}
	return r, nil
}

func ShutdownHTTPServer(server *http.Server) {
	log.Info().Msg("HTTP server shutting down.")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error occurred on http server shutting down.")
		return
	}
	log.Info().Msg("HTTP server has successfully shut down.")
}

func ShutdownGRPCServer(server *grpc.Server) {
	log.Info().Msg("gRPC server shutting down.")
	server.GracefulStop()

	log.Info().Msg("gRPC server has successfully shut down.")
}

func ClosePostgreSQL(client *sqlx.DB) {
	log.Info().Msg("Disconnecting from PostgreSQL.")
	if err := client.Close(); err != nil {
		log.Error().Err(err).Msg("Error occurred when disconnecting from PostgreSQL.")
		return
	}
	log.Info().Msg("Successfully disconnected from PostgreSQL.")
}

func CloseRedisDB(client *redisdriver.Client) {
	log.Info().Msg("Disconnecting from RedisDB.")
	if err := client.Close(); err != nil {
		log.Error().Err(err).Msg("Error occurred when disconnecting from RedisDB.")
		return
	}
	log.Info().Msg("Successfully disconnected from RedisDB.")
}
