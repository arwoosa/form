package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository provides basic MongoDB operations
type MongoRepository struct {
	client   *mongo.Client
	database string
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page     int
	PageSize int
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(client *mongo.Client, database string) *MongoRepository {
	return &MongoRepository{
		client:   client,
		database: database,
	}
}

// GetCollection returns a MongoDB collection
func (r *MongoRepository) GetCollection(name string) *mongo.Collection {
	return r.client.Database(r.database).Collection(name)
}

// Save saves a document to the specified collection
func (r *MongoRepository) Save(ctx context.Context, collection string, document interface{}) error {
	coll := r.GetCollection(collection)
	_, err := coll.InsertOne(ctx, document)
	return err
}

// FindOne finds a single document by filter
func (r *MongoRepository) FindOne(ctx context.Context, collection string, filter map[string]interface{}, result interface{}) error {
	coll := r.GetCollection(collection)
	return coll.FindOne(ctx, filter).Decode(result)
}

// FindWithPagination finds documents with pagination
func (r *MongoRepository) FindWithPagination(ctx context.Context, collection string, filter map[string]interface{}, results interface{}, pagination *PaginationOptions) (int64, error) {
	// Calculate skip based on pagination
	skip := int64(0)
	if pagination.Page > 1 {
		skip = int64((pagination.Page - 1) * pagination.PageSize)
	}

	coll := r.GetCollection(collection)

	// Get total count
	totalCount, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Find with pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pagination.PageSize)).
		SetSort(map[string]interface{}{"created_at": -1}) // Sort by created_at descending

	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return 0, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, results)
	if err != nil {
		return 0, fmt.Errorf("failed to decode documents: %w", err)
	}

	return totalCount, nil
}

// UpdateOne updates a single document
func (r *MongoRepository) UpdateOne(ctx context.Context, collection string, filter map[string]interface{}, update interface{}) error {
	coll := r.GetCollection(collection)
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

// DeleteOne deletes a single document
func (r *MongoRepository) DeleteOne(ctx context.Context, collection string, filter map[string]interface{}) error {
	coll := r.GetCollection(collection)
	_, err := coll.DeleteOne(ctx, filter)
	return err
}

// Count counts documents matching the filter
func (r *MongoRepository) Count(ctx context.Context, collection string, filter map[string]interface{}) (int64, error) {
	coll := r.GetCollection(collection)
	return coll.CountDocuments(ctx, filter)
}

// Find finds documents without pagination
func (r *MongoRepository) Find(ctx context.Context, collection string, filter map[string]interface{}, results interface{}, opts *options.FindOptions) error {
	coll := r.GetCollection(collection)
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}
