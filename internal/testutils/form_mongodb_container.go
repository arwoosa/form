package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/arwoosa/form-service/internal/models"
)

// MongoContainer wraps testcontainers MongoDB container for form service
type MongoContainer struct {
	Container testcontainers.Container
	URI       string
	Client    *mongo.Client
	Database  *mongo.Database
}

// SetupMongoContainer creates and starts a MongoDB test container
func SetupMongoContainer(ctx context.Context, t *testing.T) *MongoContainer {
	t.Helper()

	// Create MongoDB container
	mongodbContainer, err := mongodb.Run(ctx, "mongo:7",
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").
				WithOccurrence(1).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}

	// Get connection string
	uri, err := mongodbContainer.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("Failed to get MongoDB URI: %v", err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	database := client.Database("testdb")

	return &MongoContainer{
		Container: mongodbContainer,
		URI:       uri,
		Client:    client,
		Database:  database,
	}
}

// Cleanup closes the MongoDB connection and terminates the container
func (mc *MongoContainer) Cleanup(ctx context.Context, t *testing.T) {
	t.Helper()

	if mc.Client != nil {
		if err := mc.Client.Disconnect(ctx); err != nil {
			t.Logf("Failed to disconnect MongoDB client: %v", err)
		}
	}

	if mc.Container != nil {
		if err := mc.Container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate MongoDB container: %v", err)
		}
	}
}

// CleanCollections drops all collections in the test database
func (mc *MongoContainer) CleanCollections(ctx context.Context, t *testing.T) {
	t.Helper()

	collections, err := mc.Database.ListCollectionNames(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to list collections: %v", err)
	}

	for _, collection := range collections {
		if err := mc.Database.Collection(collection).Drop(ctx); err != nil {
			t.Fatalf("Failed to drop collection %s: %v", collection, err)
		}
	}
}

// CreateIndexes creates the necessary indexes for form service testing
func (mc *MongoContainer) CreateIndexes(ctx context.Context, t *testing.T) {
	t.Helper()

	// Form templates collection indexes
	templateCollection := mc.Database.Collection((&models.FormTemplate{}).TableName())

	// Merchant ID index
	_, err := templateCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"merchant_id": 1},
	})
	if err != nil {
		t.Fatalf("Failed to create merchant_id index on form_templates: %v", err)
	}

	// Compound index for merchant_id and name (for uniqueness within merchant)
	_, err = templateCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"merchant_id": 1,
			"name":        1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create compound merchant_id+name index on form_templates: %v", err)
	}

	// Forms collection indexes
	formCollection := mc.Database.Collection((&models.Form{}).TableName())

	// Merchant ID index
	_, err = formCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"merchant_id": 1},
	})
	if err != nil {
		t.Fatalf("Failed to create merchant_id index on forms: %v", err)
	}

	// Template ID index (sparse since not all forms have templates)
	_, err = formCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"template_id": 1},
		Options: options.Index().SetSparse(true),
	})
	if err != nil {
		t.Fatalf("Failed to create template_id index on forms: %v", err)
	}

	// Event ID index (sparse since not all forms have events)
	_, err = formCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"event_id": 1},
		Options: options.Index().SetSparse(true),
	})
	if err != nil {
		t.Fatalf("Failed to create event_id index on forms: %v", err)
	}

	// Compound index for merchant_id and template_id
	_, err = formCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{
			"merchant_id": 1,
			"template_id": 1,
		},
		Options: options.Index().SetSparse(true),
	})
	if err != nil {
		t.Fatalf("Failed to create compound merchant_id+template_id index on forms: %v", err)
	}
}

// GetFormTemplateCollection returns the form templates collection
func (mc *MongoContainer) GetFormTemplateCollection() *mongo.Collection {
	return mc.Database.Collection((&models.FormTemplate{}).TableName())
}

// GetFormCollection returns the forms collection
func (mc *MongoContainer) GetFormCollection() *mongo.Collection {
	return mc.Database.Collection((&models.Form{}).TableName())
}
