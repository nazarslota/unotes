package main

import (
	"context"
	"io"
	"net"
	"os"
	"time"

	"github.com/nazarslota/unotes/auth/pkg/logger"
	"github.com/nazarslota/unotes/auth/pkg/utils"
	"github.com/nazarslota/unotes/note/internal/config"
	"github.com/nazarslota/unotes/note/internal/handler"
	"github.com/nazarslota/unotes/note/internal/service"
	"github.com/nazarslota/unotes/note/internal/storage"
	"github.com/nazarslota/unotes/note/internal/storage/mongo"
)

var log logger.Logger

func init() {
	var file io.Writer
	file, err := os.OpenFile(config.C().Note.Log, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		file = io.Discard
	}
	out := io.MultiWriter(file, logger.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen})
	log = logger.NewLogger(out).With().Timestamp().Logger()
}

func main() {
	log.Info("Attempting to establish a connection with a MongoDB...")
	database, err := mongo.NewMongoDB(context.Background(), mongo.Config{
		Host:     config.C().MongoDB.Host,
		Port:     config.C().MongoDB.Port,
		Username: config.C().MongoDB.Username,
		Password: config.C().MongoDB.Password,
		Database: config.C().MongoDB.Database,
	})
	if err != nil {
		log.FatalFields("Failed to establish a connection with MongoDB.", map[string]any{"error": err})
	} else {
		log.Info("The connection to the database was successfully established.")
	}

	repositories := storage.NewRepositoryProvider(
		storage.WithMongoNoteRepository(database),
	)

	services := service.NewServices(
		service.JWTServiceOptions{AccessTokenSecret: config.C().Note.AccessTokenSecret},
		service.NoteServiceOptions{NoteRepository: repositories.NoteRepository},
	)

	grpcServerAddr := net.JoinHostPort(
		config.C().Note.HostGRPC,
		config.C().Note.PortGRPC,
	)
	restServerAddr := net.JoinHostPort(
		config.C().Note.HostREST,
		config.C().Note.PortREST,
	)

	server := handler.NewHandler(
		handler.WithServices(services),
		handler.WithGRPCServerAddr(grpcServerAddr),
		handler.WithRESTServerAddr(restServerAddr),
		handler.WithLogger(log),
	).Server()

	log.Info("Starting a gRPC server...")
	go func() {
		if err := server.ServeGRPC(); err != nil {
			log.FatalFields("Error occurred while running gRPC server.", map[string]any{"error": err})
		}
	}()

	time.Sleep(time.Second)
	log.InfoFields("The gRPC server is successfully started.", map[string]any{"address": grpcServerAddr})

	log.Info("Starting a REST server...")
	go func() {
		if err := server.ServeREST(); err != nil {
			log.FatalFields("Error occurred while running REST server.", map[string]any{"error": err})
		}
	}()

	time.Sleep(time.Second)
	log.InfoFields("The REST server is successfully started.", map[string]any{"address": restServerAddr})

	<-utils.GracefulShutdown()

	log.Info("Shutdown of the gRPC server...")
	if err := server.ShutdownGRPC(context.Background()); err != nil {
		log.ErrorFields("Error during gRPC server shutdown.", map[string]any{"error": err})
	} else {
		log.Info("gRPC server was successfully shut down.")
	}

	log.Info("Shutdown of the REST server...")
	if err := server.ShutdownREST(context.Background()); err != nil {
		log.ErrorFields("Error during REST server shutdown.", map[string]any{"error": err})
	} else {
		log.Info("REST server was successfully shut down.")
	}

	log.Info("Disconnecting from MongoDB...")
	if err := database.Client().Disconnect(context.Background()); err != nil {
		log.ErrorFields("Error during disconnecting for MongoDB.", map[string]any{"error": err})
	} else {
		log.Info("Successfully disconnected from MongoDB.")
	}

	log.Info("Shutdown completed successfully.")
}
