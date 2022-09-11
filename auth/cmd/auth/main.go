package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/udholdenhed/unotes/auth/internal/config"
	"github.com/udholdenhed/unotes/auth/internal/handler"
	"github.com/udholdenhed/unotes/auth/internal/service"
	"github.com/udholdenhed/unotes/auth/internal/storage"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	repos := storage.NewRepositoryProvider(
		storage.WithMemoryUserRepository(),
		storage.WithMemoryRefreshTokenRepository(),
	)

	services := service.NewService(&service.OAuth2ServiceOptions{
		AccessTokenSecret:      config.C().Auth.AccessTokenSecret,
		RefreshTokenSecret:     config.C().Auth.RefreshTokenSecret,
		AccessTokenExpiresIn:   config.C().Auth.AccessTokenExpiresIn,
		RefreshTokenExpiresIn:  config.C().Auth.RefreshTokenExpiresIn,
		UserRepository:         repos.UserRepository,
		RefreshTokenRepository: repos.RefreshTokenRepository,
	})

	handlers := handler.NewHandler(services, &log.Logger)

	addr := net.JoinHostPort(
		config.C().Auth.Host,
		config.C().Auth.Port,
	)

	log.Info().Msgf("Auth service starting on addr: '%s'.", addr)

	h := handlers.InitRoutes()
	server := &http.Server{
		Addr:           addr,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Error occurred while running HTTP server.")
		}
	}()

	log.Info().Msg("Auth service started.")

	<-utils.GracefulShutdown()

	log.Info().Msg("Auth service shutting down.")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Msg("Error occurred on server shutting down.")
	}

	log.Info().Msg("Auth service shutdown completed.")
}
