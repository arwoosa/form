package mgo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

func Save[T DocInter](ctx context.Context, doc T) (T, error) {
	var zero T
	if dataStore == nil {
		return zero, ErrNotConnected
	}
	newDoc, err := dataStore.Save(ctx, doc)
	if err != nil {
		return zero, fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}
	result, ok := newDoc.(T)
	if !ok {
		return zero, fmt.Errorf("%w: failed to cast to %T", ErrWriteFailed, doc)
	}
	return result, nil
}

func (m *mongoStore) Save(ctx context.Context, doc DocInter) (DocInter, error) {
	// 1. Restore the nil check for robustness.
	if v := reflect.ValueOf(doc); v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, fmt.Errorf("%w: %w", ErrInvalidDocument, errors.New("document cannot be nil"))
	}

	// 2. Restore the validation check.
	if err := doc.Validate(); err != nil {
		return doc, fmt.Errorf("%w: %w", ErrInvalidDocument, err)
	}

	// 3. Perform the database operation.
	c := m.getCollection(doc.C())
	result, err := c.InsertOne(ctx, doc)
	if err != nil {
		return doc, fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}
	doc.SetId(result.InsertedID)
	return doc, nil
}
