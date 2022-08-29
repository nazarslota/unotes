package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	redisdriver "github.com/go-redis/redis/v9"
	"github.com/udholdenhed/unotes/auth/internal/application"
	"github.com/udholdenhed/unotes/auth/internal/config"
	"github.com/udholdenhed/unotes/auth/internal/inputports"
	inputportshttp "github.com/udholdenhed/unotes/auth/internal/inputports/http"
	"github.com/udholdenhed/unotes/auth/internal/interfaceadapters/storage"
	"github.com/udholdenhed/unotes/auth/pkg/logger"
	"github.com/udholdenhed/unotes/auth/pkg/mongo"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	InitConfig()
	InitLogger()

	address := net.JoinHostPort(config.C().Host, config.C().Port)
	logger.L().Infof("trying to start server on address: %s", address)

	shutdown := make(chan struct{})
	go func() {
		<-utils.GracefulShutdown()
		close(shutdown)
	}()

	so := inputportshttp.ServerOptions{
		Address: address,
		Debug:   config.C().Debug,
	}
	RunHTTPServer(so, shutdown)
}

func InitConfig() {
	_ = config.C()
}

func InitLogger() {
	filename := config.C().LoggerOutputFile
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logger.L().SetOutput(file)
}

func RunHTTPServer(ho inputportshttp.ServerOptions, shutdown <-chan struct{}) {
	hs := NewInputPortsServices(ho, shutdown).HTTPServer
	go func() {
		if err := hs.Run(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.L().Fatalf("server shutdown error: %v", err)
		}
	}()

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := hs.Shutdown(ctx); err != nil {
		logger.L().Fatalf("server shutdown error: %v", err)
	}
	logger.L().Infof("shutdown completed successfully")
}

func NewInputPortsServices(ho inputportshttp.ServerOptions, shutdown <-chan struct{}) *inputports.Services {
	as := NewApplicationServices(shutdown)
	return inputports.NewServices(as, ho)
}

func NewApplicationServices(shutdown <-chan struct{}) *application.Services {
	rp := NewRepositoryProvider(shutdown)
	return application.NewServices(application.OAuth2ServiceOptions{
		AccessTokenSecret:      config.C().AccessTokenSecret,
		RefreshTokenSecret:     config.C().RefreshTokenSecret,
		AccessTokenExpiresIn:   config.C().AccessTokenExpiresIn,
		RefreshTokenExpiresIn:  config.C().RefreshTokenExpiresIn,
		UserRepository:         rp.UserRepository,
		RefreshTokenRepository: rp.RefreshTokenRepository,
	})
}

func NewRepositoryProvider(shutdown <-chan struct{}) *storage.RepositoryProvider {
	return storage.NewRepositoryProvider(
		storage.WithMongoUserRepository(func() *mongodriver.Collection {
			c := config.C().MongoDB
			uri := mongo.BuildURI(c.Host, c.Port, c.Username, c.Password)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			logger.L().Infof("connecting to MongoDB...")
			client, err := mongo.NewClient(ctx, uri)
			if err != nil {
				logger.L().Fatalf("filed to connect to MongoDB due to error: %v", err)
			}
			logger.L().Infof("connected to MongoDB")

			go func() {
				<-shutdown
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				logger.L().Infof("disconnecting from MongoDB database...")
				if err := client.Disconnect(ctx); err != nil {
					logger.L().Fatalf("filed to disconnect from MongoDB due to error: %v", err)
				}
				logger.L().Infof("disconnect from MongoDB")
			}()
			return client.Database(c.Database).Collection(c.UsersCollection)
		}()),
		storage.WithRedisRefreshTokenRepository(func() *redisdriver.Client {
			c := config.C().Redis
			client := redisdriver.NewClient(&redisdriver.Options{
				Addr:     c.Addr,
				Password: c.Password,
				DB:       c.Database,
			})

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			logger.L().Infof("connecting to Redis...")
			if err := client.Ping(ctx).Err(); err != nil {
				logger.L().Fatalf("filed to disconnect from MongoDB due to error: %v", err)
			}
			logger.L().Infof("connected to Redis")

			go func() {
				<-shutdown
				logger.L().Infof("disconnecting from Redis database...")
				if err := client.Close(); err != nil {
					logger.L().Fatalf("filed to colose redis connection due to error: %v", err)
				}
			}()
			return client
		}()),
	)
}
