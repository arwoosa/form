package repository

import (
	"context"
	"fmt"

	"github.com/arwoosa/form-service/pkg/vulpes/db/mgo"
)

// MongoRepository provides basic MongoDB operations
type MongoRepository struct {
	datastore mgo.Datastore
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page     int
	PageSize int
}

// NewMongoRepository creates a new MongoDB repository
func NewMongoRepository(datastore mgo.Datastore) *MongoRepository {
	return &MongoRepository{
		datastore: datastore,
	}
}

// Save saves a document to the specified collection
func (r *MongoRepository) Save(ctx context.Context, collection string, document interface{}) error {
	return r.datastore.Save(ctx, collection, document)
}

// FindOne finds a single document by filter
func (r *MongoRepository) FindOne(ctx context.Context, collection string, filter map[string]interface{}, result interface{}) error {
	return r.datastore.FindOne(ctx, collection, filter, result)
}

// FindWithPagination finds documents with pagination
func (r *MongoRepository) FindWithPagination(ctx context.Context, collection string, filter map[string]interface{}, results interface{}, pagination *PaginationOptions) (int64, error) {
	// Calculate skip based on pagination
	skip := 0
	if pagination.Page > 1 {
		skip = (pagination.Page - 1) * pagination.PageSize
	}

	// Get total count
	totalCount, err := r.datastore.Count(ctx, collection, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Find with pagination
	findOptions := &mgo.FindOptions{
		Skip:  skip,
		Limit: pagination.PageSize,
		Sort:  map[string]int{"created_at": -1}, // Sort by created_at descending
	}

	err = r.datastore.Find(ctx, collection, filter, results, findOptions)
	if err != nil {
		return 0, fmt.Errorf("failed to find documents: %w", err)
	}

	return totalCount, nil
}

// UpdateOne updates a single document
func (r *MongoRepository) UpdateOne(ctx context.Context, collection string, filter map[string]interface{}, update interface{}) error {
	return r.datastore.UpdateOne(ctx, collection, filter, update)
}

// DeleteOne deletes a single document
func (r *MongoRepository) DeleteOne(ctx context.Context, collection string, filter map[string]interface{}) error {
	return r.datastore.DeleteOne(ctx, collection, filter)
}

// Count counts documents matching the filter
func (r *MongoRepository) Count(ctx context.Context, collection string, filter map[string]interface{}) (int64, error) {
	return r.datastore.Count(ctx, collection, filter)
}

// Find finds documents without pagination
func (r *MongoRepository) Find(ctx context.Context, collection string, filter map[string]interface{}, results interface{}, options *mgo.FindOptions) error {
	return r.datastore.Find(ctx, collection, filter, results, options)
}