package mongodb

import (
	"context"
	"fmt"

	"github.com/arwoosa/form-service/conf"

	"github.com/arwoosa/vulpes/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Migration defines the structure for a collection migration
type Migration struct {
	Collection string
	Indexes    []mongo.IndexModel
}

// collection相關的index
var migrations = []Migration{
	{
		Collection: "form_templates",
		Indexes: []mongo.IndexModel{
			// Basic query index for merchant isolation and sorting
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "created_at", Value: -1},
				},
			},
			// Merchant isolation with name search
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "name", Value: 1},
				},
			},
			// Sorting by updated time for listing
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "updated_at", Value: -1},
				},
			},
			// Text search index for template names
			{
				Keys: bson.D{{Key: "name", Value: "text"}},
			},
			// Index for counting templates per merchant (business rules)
			{
				Keys: bson.D{{Key: "merchant_id", Value: 1}},
			},
		},
	},
	{
		Collection: "forms",
		Indexes: []mongo.IndexModel{
			// Basic query index for merchant isolation and sorting
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "created_at", Value: -1},
				},
			},
			// Event-based form queries
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "event_id", Value: 1},
				},
			},
			// Template-based form queries
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "template_id", Value: 1},
				},
			},
			// Forms by event for specific event queries
			{
				Keys: bson.D{{Key: "event_id", Value: 1}},
			},
			// Forms by template for template usage tracking
			{
				Keys: bson.D{{Key: "template_id", Value: 1}},
			},
			// Sorting by updated time
			{
				Keys: bson.D{
					{Key: "merchant_id", Value: 1},
					{Key: "updated_at", Value: -1},
				},
			},
			// Text search index for form names
			{
				Keys: bson.D{{Key: "name", Value: "text"}},
			},
		},
	},
}

// Migrate runs all the defined migrations.
func Migrate(client *mongo.Client, cfg *conf.MongodbConfig) error {
	log.Info("Running MongoDB migrations...")
	db := client.Database(cfg.DB)

	for _, m := range migrations {
		coll := db.Collection(m.Collection)

		if len(m.Indexes) > 0 {
			_, err := coll.Indexes().CreateMany(context.Background(), m.Indexes)
			if err != nil {
				return fmt.Errorf("failed to create indexes for collection '%s': %w", m.Collection, err)
			}
		}
	}

	return nil
}
