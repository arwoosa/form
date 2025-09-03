// Package mgo provides a high-level abstraction layer over the official MongoDB Go driver,
// simplifying connection management, document operations, and schema definitions.
package mgo

// DocInter defines a standard interface for all database documents.
// By embedding the Index interface, it requires all documents to specify
// their collection name, index definitions, and validation logic.
// This promotes a consistent and self-describing data model.
type DocInter interface {
	Index
	// Validate performs business logic validation on the document's fields.
	Validate() error
	GetId() any
	SetId(any)
}
