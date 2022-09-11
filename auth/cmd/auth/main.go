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
	"github.com/udholdenhed/unotes/auth/internal/storage/mongo"
	"github.com/udholdenhed/unotes/auth/internal/storage/redis"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Connecting to MongoDB.")
	mongoDBClient, err := mongo.NewMongoDBClient(context.Background(), &mongo.Config{
		Host:     config.C().MongoDB.Host,
		Port:     config.C().MongoDB.Port,
		Username: config.C().MongoDB.Username,
		Password: config.C().MongoDB.Password,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Filed to connect to MongoDB.")
	}
	log.Info().Msg("Connected to MongoDB.")

	log.Info().Msg("Connecting to Redis.")
	redisClient, err := redis.NewRedisClient(context.Background(), &redis.Config{
		Addr:     "",
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Filed to connect to Redis.")
	}
	log.Info().Msg("Connected to Redis.")

	repos := storage.NewRepositoryProvider(
		storage.WithMongoDBUserRepository(func() *mongodriver.Collection {
			return mongoDBClient.Database(config.C().MongoDB.Database).Collection("users")
		}()),
		storage.WithRedisRefreshTokenRepository(redisClient),
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

	log.Info().Msgf("Starting HTTP server on addr: %q...", addr)

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

	log.Info().Msg("Server started.")

	<-utils.GracefulShutdown()

	log.Info().Msg("Shutting down...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error occurred on server shutting down.")
	}

	log.Info().Msg("Disconnecting from MongoDB...")
	if err := mongoDBClient.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("Error occurred when disconnecting from the MongoDB.")
	} else {
		log.Info().Msg("Disconnected from MongoDB.")
	}

	log.Info().Msg("Disconnecting from Redis...")
	if err := redisClient.Close(); err != nil {
		log.Error().Err(err).Msg("Error occurred when disconnecting from the Redis.")
	} else {
		log.Info().Msg("Disconnected from Redis.")
	}

	log.Info().Msg("Shutdown completed.")
}
