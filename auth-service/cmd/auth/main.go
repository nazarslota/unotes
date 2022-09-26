package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	redisdriver "github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/udholdenhed/unotes/auth-service/internal/config"
	"github.com/udholdenhed/unotes/auth-service/internal/handler"
	"github.com/udholdenhed/unotes/auth-service/internal/service"
	"github.com/udholdenhed/unotes/auth-service/internal/storage"
	"github.com/udholdenhed/unotes/auth-service/internal/storage/postgres"
	"github.com/udholdenhed/unotes/auth-service/internal/storage/redis"
	"github.com/udholdenhed/unotes/auth-service/pkg/utils"
)

func main() {
	initLogger()

	postgresDB, err := newPostgreSQL()
	if err != nil {
		log.Fatal().Err(err).Msg("Filed to connect to PostgreSQL.")
	}

	redisDB, err := newRedis()
	if err != nil {
		log.Fatal().Err(err).Msg("Filed to connect to Redis.")
	}

	repos := storage.NewRepositoryProvider(
		storage.WithPostgreSQLUserRepository(postgresDB),
		storage.WithRedisRefreshTokenRepository(redisDB),
	)

	services := service.NewService(&service.OAuth2ServiceOptions{
		AccessTokenSecret:      config.C().Auth.AccessTokenSecret,
		RefreshTokenSecret:     config.C().Auth.RefreshTokenSecret,
		AccessTokenExpiresIn:   config.C().Auth.AccessTokenExpiresIn,
		RefreshTokenExpiresIn:  config.C().Auth.RefreshTokenExpiresIn,
		UserRepository:         repos.UserRepository,
		RefreshTokenRepository: repos.RefreshTokenRepository,
	})

	server := &http.Server{
		Addr:           net.JoinHostPort(config.C().Auth.Host, config.C().Auth.Port),
		Handler:        handler.NewHandler(services, &log.Logger).H(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes, // 1 MB
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			switch {
			case errors.Is(err, http.ErrServerClosed):
			default:
				log.Fatal().Err(err).Msg("Error occurred while running HTTP server.")
			}
		}
	}()

	log.Info().Msg("Server started.")

	<-utils.GracefulShutdown()
	log.Info().Msg("Shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error occurred on server shutting down.")
	}

	if err := postgresDB.Close(); err != nil {
		log.Error().Err(err).Msg("Error occurred when disconnecting from the PostgreSQL.")
	}

	if err := redisDB.Close(); err != nil {
		log.Error().Err(err).Msg("Error occurred when disconnecting from the Redis.")
	}

	log.Info().Msg("Shutdown completed.")
}

func initLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC1123,
	})
}

func newPostgreSQL() (*sqlx.DB, error) {
	return postgres.NewPostgreSQL(context.Background(), &postgres.Config{
		Host:     config.C().PostgreSQL.Host,
		Port:     config.C().PostgreSQL.Port,
		Username: config.C().PostgreSQL.Username,
		Password: config.C().PostgreSQL.Password,
		DBName:   config.C().PostgreSQL.DBName,
		SSLMode:  config.C().PostgreSQL.SSLMode,
	})
}

func newRedis() (*redisdriver.Client, error) {
	return redis.NewRedis(context.Background(), &redis.Config{
		Addr:     config.C().Redis.Addr,
		Password: config.C().Redis.Password,
		DB:       config.C().Redis.DB,
	})
}
