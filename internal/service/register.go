package service

import (
	"google.golang.org/grpc"

	"github.com/arwoosa/form-service/conf"
	pb "github.com/arwoosa/form-service/gen/pb/form"
	"github.com/arwoosa/form-service/internal/dao/mongodb"
	"github.com/arwoosa/form-service/internal/dao/repository"

	"github.com/arwoosa/vulpes/ezgrpc"
	"github.com/arwoosa/vulpes/log"
)

// This file registers the services with the Vulpes framework

// RegisterFormServices registers form services
func RegisterFormServices(appConfig *conf.AppConfig) {
	// Register form gRPC services
	ezgrpc.InjectGrpcService(func(s grpc.ServiceRegistrar) {
		registerFormServices(s, appConfig)
	})

	// Register form gRPC-Gateway handlers
	ezgrpc.RegisterHandlerFromEndpoint(pb.RegisterFormServiceHandlerFromEndpoint)
}

// registerFormServices sets up and registers form related gRPC services
func registerFormServices(s grpc.ServiceRegistrar, appConfig *conf.AppConfig) {
	if appConfig == nil {
		log.Warn("Form services initialized with nil config - using mock services")
		grpcServer := NewGRPCFormServer(nil, nil)
		pb.RegisterFormServiceServer(s, grpcServer)
		return
	}

	// Get MongoDB singleton (should be initialized by main)
	mongoClient := mongodb.GetMongoDB()
	if mongoClient == nil {
		log.Warn("Form services initialized without MongoDB - using mock services")
		grpcServer := NewGRPCFormServer(nil, nil)
		pb.RegisterFormServiceServer(s, grpcServer)
		return
	}

	log.Info("Form services initialized with MongoDB connection")

	// Initialize repositories
	mongoRepo := repository.NewMongoRepository(mongoClient, appConfig.DB)
	templateRepo := repository.NewFormTemplateRepository(mongoRepo)
	formRepo := repository.NewFormRepository(mongoRepo)

	// Initialize services
	templateService := NewFormTemplateService(templateRepo, appConfig)
	formService := NewFormService(formRepo, templateRepo, appConfig)

	// Create gRPC server with the services
	grpcServer := NewGRPCFormServer(templateService, formService)

	// Register form service
	pb.RegisterFormServiceServer(s, grpcServer)
}
