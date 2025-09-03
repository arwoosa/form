package mgo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Find[T DocInter](ctx context.Context, doc T, filter any, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	if dataStore == nil {
		return nil, ErrNotConnected
	}
	result, err := dataStore.Find(ctx, doc.C(), filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	var ret []T
	err = result.All(ctx, &ret)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return ret, nil
}

func FindOne[T DocInter](ctx context.Context, doc T, filter any, opts ...options.Lister[options.FindOneOptions]) error {
	if dataStore == nil {
		return ErrNotConnected
	}
	err := dataStore.FindOne(ctx, doc.C(), filter, opts...).Decode(&doc)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return nil
}

func FindById[T DocInter](ctx context.Context, doc T) error {
	if dataStore == nil {
		return ErrNotConnected
	}
	return FindOne(ctx, doc, bson.M{"_id": doc.GetId()})
}

func (m *mongoStore) Find(ctx context.Context, collectionName string, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	collection := m.getCollection(collectionName)
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFailed, err)
	}
	return cursor, nil
}

func (m *mongoStore) FindOne(ctx context.Context, collectionName string, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	collection := m.getCollection(collectionName)
	return collection.FindOne(ctx, filter, opts...)
}
