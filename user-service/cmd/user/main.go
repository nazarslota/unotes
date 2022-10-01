package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/udholdenhed/unotes/user-service/internal/config"
	"github.com/udholdenhed/unotes/user-service/internal/handler"
	"github.com/udholdenhed/unotes/user-service/internal/service"
	"github.com/udholdenhed/unotes/user-service/internal/storage"
	"github.com/udholdenhed/unotes/user-service/internal/storage/postgres"
	"github.com/udholdenhed/unotes/user-service/pkg/utils"
)

func main() {
	initLogger()

	postgreSQL, err := newPostgreSQL()
	if err != nil {
		log.Fatal().Err(err).Msg("Filed to connect to PostgreSQL.")
	}

	repos := storage.NewRepositoryProvider(
		storage.WithPostgreSQLUserRepository(postgreSQL),
	)

	services := service.NewService(&service.UserServiceOptions{
		UserRepository: repos.UserRepository,
	})

	server := &http.Server{
		Addr:           net.JoinHostPort(config.C().UserService.Host, config.C().UserService.Port), // ":8082",
		Handler:        handler.NewHandler(services, &log.Logger).H(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
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
