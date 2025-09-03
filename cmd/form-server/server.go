package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/arwoosa/form-service/internal/dao/mongodb"
	"github.com/arwoosa/form-service/internal/service"
	"github.com/arwoosa/vulpes/ezgrpc"
	"github.com/arwoosa/vulpes/log"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the form service server",
	Long:  `Start the form service server with both gRPC and HTTP APIs`,
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	log.Info("Starting Form microservice")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appConfig := GetAppConfig()

	log.Info("Starting form service server",
		log.String("version", appConfig.Version),
		log.String("mode", appConfig.Mode),
		log.Int("port", appConfig.Port))

	// Initialize MongoDB singleton first - Form service requires database
	if _, err := mongodb.InitMongoDB(ctx, appConfig.MongodbConfig); err != nil {
		log.Error("Failed to initialize MongoDB", log.Err(err))
		log.Fatal("Form service requires MongoDB connection - cannot start without database")
	}

	// Register services
	service.RegisterFormServices(appConfig)

	ezgrpc.SetServeMuxOpts(
		ezgrpc.DefaultHeaderMatcher,
		ezgrpc.OutgoingHeaderMatcher,
	)

	// Channel to listen for server errors
	errChan := make(chan error, 1)

	// Run the gRPC + Gateway server in a goroutine
	go func() {
		if err := ezgrpc.RunGrpcGateway(ctx, appConfig.Port); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case <-ctx.Done():
		log.Info("Shutdown signal received, shutting down server gracefully...")
	case err := <-errChan:
		log.Fatal("failed to run form server", log.Err(err))
	}

	log.Info("Server shut down gracefully")
}
