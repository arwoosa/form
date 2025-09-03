package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/arwoosa/form-service/internal/dao/mongodb"
	"github.com/arwoosa/form-service/internal/service"

	"github.com/arwoosa/vulpes/ezgrpc"
	vulpeslog "github.com/arwoosa/vulpes/log"
)

// publicCmd represents the public command
var publicCmd = &cobra.Command{
	Use:   "public",
	Short: "Start the public API server",
	Long: `Start the public API server that provides read-only access to published events.

This includes search and retrieval operations for public events, typically used by
web applications and mobile apps for end users.

API endpoints:
- /events/* (public read-only API)

This service is intended for public access and provides only published events.`,
	Run: runPublicServer,
}

func runPublicServer(cmd *cobra.Command, args []string) {
	vulpeslog.Info("Starting Event microservice - Public Mode")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appConfig := GetAppConfig()

	// Initialize MongoDB singleton first - Public service can fallback to mock for graceful degradation
	if _, err := mongodb.InitMongoDB(ctx, appConfig.MongodbConfig); err != nil {
		vulpeslog.Error("Failed to initialize MongoDB", vulpeslog.Err(err))
		vulpeslog.Fatal("Public service requires MongoDB connection - cannot start without database")
	}

	// Register only public services
	service.RegisterPublicServices(appConfig)

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
		vulpeslog.Info("Shutdown signal received, shutting down server gracefully...")
	case err := <-errChan:
		vulpeslog.Fatal("failed to run public server", vulpeslog.Err(err))
	}

	vulpeslog.Info("Server shut down gracefully")
}

func init() {
	rootCmd.AddCommand(publicCmd)
}
