package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/arwoosa/form-service/gen/pb/form"
	"github.com/arwoosa/form-service/conf"
	"github.com/arwoosa/form-service/internal/dao/mongodb"
	"github.com/arwoosa/form-service/internal/dao/repository"
	"github.com/arwoosa/form-service/internal/service"
	"github.com/arwoosa/form-service/pkg/vulpes/ezgrpc"
	"github.com/arwoosa/form-service/pkg/vulpes/log"
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
	config := GetAppConfig()
	
	log.Info("Starting form service server",
		log.String("version", config.Version),
		log.String("mode", config.Mode),
		log.Int("port", config.Port))

	// Initialize MongoDB connection
	mongoDS, err := mongodb.NewMongoDB(config)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB", log.Err(err))
	}
	defer mongoDS.Close()

	// Initialize repositories
	mongoRepo := repository.NewMongoRepository(mongoDS)
	templateRepo := repository.NewFormTemplateRepository(mongoRepo)
	formRepo := repository.NewFormRepository(mongoRepo)

	// Initialize services
	templateService := service.NewFormTemplateService(templateRepo, config)
	formService := service.NewFormService(formRepo, templateRepo, config)

	// Initialize gRPC server
	grpcServer := service.NewGRPCFormServer(templateService, formService)

	// Set up and run servers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startGRPCServer(ctx, config, grpcServer); err != nil {
			log.Error("gRPC server error", log.Err(err))
		}
	}()

	// Start HTTP Gateway server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startHTTPServer(ctx, config); err != nil {
			log.Error("HTTP server error", log.Err(err))
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down servers...")
	cancel()
	wg.Wait()
	log.Info("Servers shut down gracefully")
}

func startGRPCServer(ctx context.Context, config *conf.AppConfig, grpcServer *service.GRPCFormServer) error {
	grpcPort := config.Port + 1 // Use port+1 for gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port %d: %w", grpcPort, err)
	}

	// Create gRPC server with Vulpes interceptors
	server := ezgrpc.NewGrpcServer()

	// Register form service
	pb.RegisterFormServiceServer(server, grpcServer)

	log.Info("Starting gRPC server", log.Int("port", grpcPort))

	// Start server in a goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Error("gRPC server failed", log.Err(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Info("Shutting down gRPC server...")
	server.GracefulStop()
	return nil
}

func startHTTPServer(ctx context.Context, config *conf.AppConfig) error {
	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux(
		runtime.WithHealthEndpointAt(http.MethodGet, "/health"),
	)

	// Connect to gRPC server
	grpcPort := config.Port + 1
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterFormServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", grpcPort),
		opts,
	)
	if err != nil {
		return fmt.Errorf("failed to register gRPC-Gateway handler: %w", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	log.Info("Starting HTTP server", log.Int("port", config.Port))

	// Start server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("HTTP server failed", log.Err(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	log.Info("Shutting down HTTP server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}