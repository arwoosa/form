package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MgoAggregate interface {
	GetPipeline(q bson.M) mongo.Pipeline
	Index
}

func PipeFind[T MgoAggregate](ctx context.Context, aggr T, filter bson.M) ([]T, error) {
	if dataStore == nil {
		return nil, ErrNotConnected
	}
	sortCursor, err := dataStore.PipeFind(ctx, aggr.C(), aggr.GetPipeline(filter))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	var slice []T
	err = sortCursor.All(ctx, &slice)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return slice, nil
}

func PipeFindOne[T MgoAggregate](ctx context.Context, aggr T, filter bson.M) error {
	if dataStore == nil {
		return ErrNotConnected
	}
	err := dataStore.PipeFindOne(ctx, aggr.C(), aggr.GetPipeline(filter)).Decode(&aggr)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return nil
}

func (m *mongoStore) PipeFind(ctx context.Context, collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	c := m.getCollection(collection)
	sortCursor, err := c.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return sortCursor, nil
}

func (m *mongoStore) PipeFindOne(ctx context.Context, collection string, pipeline mongo.Pipeline) *mongo.SingleResult {
	c := m.getCollection(collection)
	sortCursor, err := c.Aggregate(ctx, pipeline)
	if err != nil {
		return mongo.NewSingleResultFromDocument(bson.D{}, err, nil)
	}
	if !sortCursor.Next(ctx) {
		return mongo.NewSingleResultFromDocument(bson.D{}, mongo.ErrNoDocuments, nil)
	}
	return mongo.NewSingleResultFromDocument(sortCursor.Current, nil, nil)
}
