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
	log "github.com/sirupsen/logrus"
	"github.com/udholdenhed/unotes/auth/internal/config"
	handlergrpc "github.com/udholdenhed/unotes/auth/internal/handler/grpc"
	handlerhttp "github.com/udholdenhed/unotes/auth/internal/handler/http"
	"github.com/udholdenhed/unotes/auth/internal/service"
	"github.com/udholdenhed/unotes/auth/internal/storage"
	"github.com/udholdenhed/unotes/auth/internal/storage/postgres"
	"github.com/udholdenhed/unotes/auth/internal/storage/redis"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
)

func main() {
	InitLogger()

	postgresDB, err := NewPostgreSQL()
	if err != nil {
		log.Info("Filed to connect to PostgreSQL.")
	}

	redisDB, err := NewRedisDB()
	if err != nil {
		log.WithError(err).Error("Filed to connect to Redis.")
	}

	repositories := storage.NewRepositoryProvider(
		storage.WithPostgreSQLUserRepository(postgresDB),
		storage.WithRedisRefreshTokenRepository(redisDB),
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
	serverHTTP := handlerhttp.NewHandler(
		handlerhttp.WithAddress(addrHTTP),
		handlerhttp.WithService(services),
		handlerhttp.WithLogger(log.StandardLogger()),
	).S()

	go func() {
		if err := serverHTTP.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Warn("Error occurred while running HTTP server.")
		}
	}()
	log.Info("HTTP server started.")

	addrGRPC := net.JoinHostPort(
		config.C().Auth.HostGRPC,
		config.C().Auth.PortGRPC,
	)
	serverGRPC := handlergrpc.NewHandler(
		handlergrpc.WithAddress(addrGRPC),
		handlergrpc.WithService(services),
		handlergrpc.WithLogger(log.StandardLogger()),
	).S()

	go func() {
		if err := serverGRPC.Serve(); err != nil {
			log.WithError(err).Warn("Error occurred while running gRPC server.")
		}
	}()
	log.Info("gRPC server started.")

	<-utils.GracefulShutdown()
	log.Info("Shutting down.")

	ShutdownHTTPServer(serverHTTP)
	ShutdownGRPCServer(serverGRPC)

	ClosePostgreSQL(postgresDB)
	CloseRedisDB(redisDB)

	log.Info("Shutdown completed.")
}

func InitLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC1123,
	})
	log.SetOutput(os.Stdout)
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
		return nil, fmt.Errorf("failed to create a postgres client: %w", err)
	}
	return p, nil
}

func ClosePostgreSQL(client *sqlx.DB) {
	log.Info("Disconnecting from PostgreSQL.")
	if err := client.Close(); err != nil {
		log.WithError(err).Warn("Error occurred when disconnecting from PostgreSQL.")
		return
	}
	log.Info("Successfully disconnected from PostgreSQL.")
}

func NewRedisDB() (*redisdriver.Client, error) {
	r, err := redis.NewRedis(context.Background(), &redis.Config{
		Addr:     config.C().Redis.Addr,
		Password: config.C().Redis.Password,
		DB:       config.C().Redis.DB,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create a redis client: %w", err)
	}
	return r, nil
}

func CloseRedisDB(client *redisdriver.Client) {
	log.Info("Disconnecting from RedisDB.")
	if err := client.Close(); err != nil {
		log.WithError(err).Warn("Error occurred when disconnecting from RedisDB.")
		return
	}
	log.Info("Successfully disconnected from RedisDB.")
}

func ShutdownHTTPServer(server *handlerhttp.Server) {
	log.Info("HTTP server shutting down.")
	if err := server.Shutdown(context.Background()); err != nil {
		log.WithError(err).Warn("Error occurred on HTTP server shutting down.")
		return
	}
	log.Info("HTTP server has successfully shut down.")
}

func ShutdownGRPCServer(server *handlergrpc.Server) {
	log.Info("gRPC server shutting down.")
	if err := server.Shutdown(context.Background()); err != nil {
		log.WithError(err).Warn("Error occurred on gRPC server shutting down.")
		return
	}
	log.Info("gRPC server has successfully shut down.")
}
