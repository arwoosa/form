package service

import (
	"google.golang.org/grpc"

	"github.com/arwoosa/form-service/conf"
	consolepb "github.com/arwoosa/form-service/gen/pb/console"
	publicpb "github.com/arwoosa/form-service/gen/pb/public"
	"github.com/arwoosa/form-service/internal/dao/mongodb"
	"github.com/arwoosa/form-service/internal/dao/repository"

	"github.com/arwoosa/vulpes/ezgrpc"
	"github.com/arwoosa/vulpes/log"
)

// This file registers the services with the Vulpes framework

// RegisterConsoleServices registers only the console (management) services
func RegisterConsoleServices(appConfig *conf.AppConfig) {
	// Register console gRPC services
	ezgrpc.InjectGrpcService(func(s grpc.ServiceRegistrar) {
		registerConsoleServices(s, appConfig)
	})

	// Register console gRPC-Gateway handlers
	ezgrpc.RegisterHandlerFromEndpoint(consolepb.RegisterEventServiceHandlerFromEndpoint)
}

// RegisterPublicServices registers only the public services
func RegisterPublicServices(appConfig *conf.AppConfig) {
	// Register public gRPC services
	ezgrpc.InjectGrpcService(func(s grpc.ServiceRegistrar) {
		registerPublicServices(s, appConfig)
	})

	// Register public gRPC-Gateway handlers
	ezgrpc.RegisterHandlerFromEndpoint(publicpb.RegisterPublicEventServiceHandlerFromEndpoint)
}

// registerConsoleServices sets up and registers only console (EventService) related gRPC services
func registerConsoleServices(s grpc.ServiceRegistrar, appConfig *conf.AppConfig) {
	if appConfig == nil {
		log.Warn("Console services initialized with nil config - using mock services")
		mockOrderService := NewMockOrderServiceClient(false, nil)
		eventSvc := &EventService{eventRepo: nil, sessionService: nil, orderService: mockOrderService}
		consolepb.RegisterEventServiceServer(s, NewEventServiceServer(eventSvc, nil))
		consolepb.RegisterInternalServiceServer(s, NewInternalServiceServer(nil, nil))
		return
	}

	// Get MongoDB singleton (should be initialized by main)
	mongoClient := mongodb.GetMongoDB()
	if mongoClient == nil {
		log.Warn("Console services initialized without MongoDB - using mock services")
		mockOrderService := NewMockOrderServiceClient(false, nil)
		eventSvc := &EventService{eventRepo: nil, sessionService: nil, orderService: mockOrderService}
		consolepb.RegisterEventServiceServer(s, NewEventServiceServer(eventSvc, appConfig.PaginationConfig))
		consolepb.RegisterInternalServiceServer(s, NewInternalServiceServer(nil, nil))
		return
	}

	// Initialize repositories
	log.Info("Console services initialized with MongoDB connection")
	eventRepo := repository.NewMongoEventRepository(mongoClient, appConfig.DB, appConfig.PaginationConfig)
	sessionRepo := repository.NewMongoSessionRepository(mongoClient, appConfig.DB)

	// Initialize external services
	var orderService OrderServiceClient
	if appConfig.ExternalConfig != nil {
		log.Info("Console services using real order service")
		orderService = NewOrderServiceClient(appConfig.OrderService)
	} else {
		log.Warn("Console services initialized without external config - using mock order service")
		orderService = NewMockOrderServiceClient(false, nil)
	}

	// Initialize business services
	sessionSvc := NewSessionService(sessionRepo, eventRepo)
	eventSvc := NewEventService(eventRepo, sessionSvc, orderService)

	// Register console services
	consolepb.RegisterEventServiceServer(s, NewEventServiceServer(eventSvc, appConfig.PaginationConfig))
	consolepb.RegisterInternalServiceServer(s, NewInternalServiceServer(eventRepo, sessionRepo))
}

// registerPublicServices sets up and registers only public (PublicEventService) related gRPC services
func registerPublicServices(s grpc.ServiceRegistrar, appConfig *conf.AppConfig) {
	if appConfig == nil {
		log.Warn("Public services initialized with nil config - using mock services")
		publicSvc := &PublicService{eventRepo: nil, sessionService: nil, paginationConfig: nil}
		publicpb.RegisterPublicEventServiceServer(s, NewPublicEventServiceServer(publicSvc))
		return
	}

	// Get MongoDB singleton (should be initialized by main)
	mongoClient := mongodb.GetMongoDB()
	if mongoClient == nil {
		log.Warn("Public services initialized without MongoDB - using mock services")
		publicSvc := &PublicService{eventRepo: nil, sessionService: nil, paginationConfig: appConfig.PaginationConfig}
		publicpb.RegisterPublicEventServiceServer(s, NewPublicEventServiceServer(publicSvc))
		return
	}

	// Initialize repositories
	log.Info("Public services initialized with MongoDB connection")
	eventRepo := repository.NewMongoEventRepository(mongoClient, appConfig.DB, appConfig.PaginationConfig)
	sessionRepo := repository.NewMongoSessionRepository(mongoClient, appConfig.DB)

	// Initialize business services
	sessionSvc := NewSessionService(sessionRepo, eventRepo)
	publicSvc := NewPublicService(eventRepo, sessionSvc, appConfig.PaginationConfig)

	// Register only PublicEventService (public API)
	publicpb.RegisterPublicEventServiceServer(s, NewPublicEventServiceServer(publicSvc))
}
