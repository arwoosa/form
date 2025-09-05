package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form/internal/models"
)

// FormTemplateRepository defines the interface for form template data access
type FormTemplateRepository interface {
	// Create a new form template
	Create(ctx context.Context, template *models.FormTemplate) error

	// Find form template by ID
	FindByID(ctx context.Context, templateID primitive.ObjectID) (*models.FormTemplate, error)

	// Find form templates with pagination by merchant ID
	FindByMerchantID(ctx context.Context, options *models.FormTemplateQueryOptions) ([]*models.FormTemplate, int64, error)

	// Update form template
	Update(ctx context.Context, template *models.FormTemplate) error

	// Delete form template by ID
	Delete(ctx context.Context, templateID primitive.ObjectID) error

	// Count templates by merchant ID (for business rule validation)
	CountByMerchantID(ctx context.Context, merchantID string) (int64, error)

	// Check if template exists by ID
	Exists(ctx context.Context, templateID primitive.ObjectID) (bool, error)

	// Duplicate a template with new name
	Duplicate(ctx context.Context, sourceID primitive.ObjectID, nameSuffix, createdBy, merchantID string) (*models.FormTemplate, error)
}

// NewFormTemplateRepository creates a new form template repository implementation
func NewFormTemplateRepository(mongoRepo *MongoRepository) FormTemplateRepository {
	return &mongoFormTemplateRepository{
		mongoRepo: mongoRepo,
	}
}

type mongoFormTemplateRepository struct {
	mongoRepo *MongoRepository
}

// Create implements FormTemplateRepository.Create
func (r *mongoFormTemplateRepository) Create(ctx context.Context, template *models.FormTemplate) error {
	now := time.Now()
	template.SetCreatedAt(now)
	template.SetUpdatedAt(now)

	if template.ID.IsZero() {
		template.ID = primitive.NewObjectID()
	}

	return r.mongoRepo.Save(ctx, template.TableName(), template)
}

// FindByID implements FormTemplateRepository.FindByID
func (r *mongoFormTemplateRepository) FindByID(ctx context.Context, templateID primitive.ObjectID) (*models.FormTemplate, error) {
	filter := map[string]interface{}{
		"_id": templateID,
	}

	var template models.FormTemplate
	err := r.mongoRepo.FindOne(ctx, template.TableName(), filter, &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// FindByMerchantID implements FormTemplateRepository.FindByMerchantID
func (r *mongoFormTemplateRepository) FindByMerchantID(ctx context.Context, options *models.FormTemplateQueryOptions) ([]*models.FormTemplate, int64, error) {
	filter := map[string]interface{}{
		"merchant_id": options.MerchantID,
	}

	var templates []*models.FormTemplate
	pagination := &PaginationOptions{
		Page:      options.Page,
		PageSize:  options.PageSize,
		SortBy:    options.SortBy,
		SortOrder: options.SortOrder,
	}

	count, err := r.mongoRepo.FindWithPagination(ctx, models.FormTemplate{}.TableName(), filter, &templates, pagination)
	if err != nil {
		return nil, 0, err
	}

	return templates, count, nil
}

// Update implements FormTemplateRepository.Update
func (r *mongoFormTemplateRepository) Update(ctx context.Context, template *models.FormTemplate) error {
	template.SetUpdatedAt(time.Now())

	filter := map[string]interface{}{
		"_id":         template.ID,
		"merchant_id": template.MerchantID,
	}

	return r.mongoRepo.UpdateOne(ctx, template.TableName(), filter, template)
}

// Delete implements FormTemplateRepository.Delete
func (r *mongoFormTemplateRepository) Delete(ctx context.Context, templateID primitive.ObjectID) error {
	filter := map[string]interface{}{
		"_id": templateID,
	}

	return r.mongoRepo.DeleteOne(ctx, models.FormTemplate{}.TableName(), filter)
}

// CountByMerchantID implements FormTemplateRepository.CountByMerchantID
func (r *mongoFormTemplateRepository) CountByMerchantID(ctx context.Context, merchantID string) (int64, error) {
	filter := map[string]interface{}{
		"merchant_id": merchantID,
	}

	return r.mongoRepo.Count(ctx, models.FormTemplate{}.TableName(), filter)
}

// Exists implements FormTemplateRepository.Exists
func (r *mongoFormTemplateRepository) Exists(ctx context.Context, templateID primitive.ObjectID) (bool, error) {
	count, err := r.mongoRepo.Count(ctx, models.FormTemplate{}.TableName(), map[string]interface{}{
		"_id": templateID,
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Duplicate implements FormTemplateRepository.Duplicate
func (r *mongoFormTemplateRepository) Duplicate(ctx context.Context, sourceID primitive.ObjectID, nameSuffix, createdBy, merchantID string) (*models.FormTemplate, error) {
	// First, find the source template
	source, err := r.FindByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	// Create a duplicate with new name and metadata
	duplicate := &models.FormTemplate{
		ID:         primitive.NewObjectID(),
		Name:       source.Name + nameSuffix,
		MerchantID: merchantID,
		Schema:     source.Schema,
		UISchema:   source.UISchema,
		CreatedBy:  createdBy,
		UpdatedBy:  createdBy,
	}

	err = r.Create(ctx, duplicate)
	if err != nil {
		return nil, err
	}

	return duplicate, nil
}
