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
	handlergrpc "github.com/nazarslota/unotes/note/internal/handler/grpc"
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
	log.Info("Attempting to establish a connection to a MongoDB...")
	database, err := mongo.NewMongoDB(context.Background(), mongo.Config{
		Host:     config.C().MongoDB.Host,
		Port:     config.C().MongoDB.Port,
		Username: config.C().MongoDB.Username,
		Password: config.C().MongoDB.Password,
		Database: config.C().MongoDB.Database,
	})
	if err != nil {
		log.ErrorFields("Attempting to establish a connection to a mongo database.", map[string]any{"error": err})
	} else {
		log.Info("The connection to the database was successfully established.")
	}

	repositories := storage.NewRepositoryProvider(
		storage.WithMongoNoteRepository(database),
	)

	services := service.NewServices(
		service.NoteServiceOptions{
			NoteRepository: repositories.NoteRepository,
		},
	)

	address := net.JoinHostPort(config.C().Note.HostGRPC, config.C().Note.PortGRPC)
	server := handlergrpc.NewHandler(
		handlergrpc.WithServices(services),
		handlergrpc.WithAddress(address),
		handlergrpc.WithLogger(log),
	).S()

	log.Info("Starting a gRPC server...")
	go func() {
		if err := server.Serve(); err != nil {
			log.FatalFields("Error occurred while running gRPC server.", map[string]any{"error": err})
		}
	}()
	log.InfoFields("The gRPC server is successfully started.", map[string]any{"address": address})

	<-utils.GracefulShutdown()

	log.Info("Shutdown of the gRPC server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.ErrorFields("Error during gRPC server shutdown.", map[string]any{"error": err})
	} else {
		log.Info("gRPC server was successfully shut down.")
	}
}
