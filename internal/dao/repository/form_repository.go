package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form-service/internal/models"
)

// FormRepository defines the interface for form data access
type FormRepository interface {
	// Create a new form
	Create(ctx context.Context, form *models.Form) error

	// Find form by ID and merchant ID
	FindByID(ctx context.Context, formID primitive.ObjectID, merchantID string) (*models.Form, error)

	// Find forms with pagination and optional filters
	Find(ctx context.Context, options *models.FormQueryOptions) ([]*models.Form, int64, error)

	// Update form
	Update(ctx context.Context, form *models.Form) error

	// Delete form by ID and merchant ID
	Delete(ctx context.Context, formID primitive.ObjectID, merchantID string) error

	// Check if form exists by ID and merchant ID
	Exists(ctx context.Context, formID primitive.ObjectID, merchantID string) (bool, error)

	// Find forms by event ID
	FindByEventID(ctx context.Context, eventID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error)

	// Find forms by template ID
	FindByTemplateID(ctx context.Context, templateID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error)

	// Count forms using a specific template (useful for template deletion validation)
	CountByTemplateID(ctx context.Context, templateID primitive.ObjectID, merchantID string) (int64, error)
}

// NewFormRepository creates a new form repository implementation
func NewFormRepository(mongoRepo *MongoRepository) FormRepository {
	return &mongoFormRepository{
		mongoRepo: mongoRepo,
	}
}

type mongoFormRepository struct {
	mongoRepo *MongoRepository
}

// Create implements FormRepository.Create
func (r *mongoFormRepository) Create(ctx context.Context, form *models.Form) error {
	now := time.Now()
	form.SetCreatedAt(now)
	form.SetUpdatedAt(now)

	if form.ID.IsZero() {
		form.ID = primitive.NewObjectID()
	}

	return r.mongoRepo.Save(ctx, form.TableName(), form)
}

// FindByID implements FormRepository.FindByID
func (r *mongoFormRepository) FindByID(ctx context.Context, formID primitive.ObjectID, merchantID string) (*models.Form, error) {
	filter := map[string]interface{}{
		"_id":         formID,
		"merchant_id": merchantID,
	}

	var form models.Form
	err := r.mongoRepo.FindOne(ctx, form.TableName(), filter, &form)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

// Find implements FormRepository.Find
func (r *mongoFormRepository) Find(ctx context.Context, options *models.FormQueryOptions) ([]*models.Form, int64, error) {
	filter := map[string]interface{}{
		"merchant_id": options.MerchantID,
	}

	// Add optional filters
	if options.EventID != nil && !options.EventID.IsZero() {
		filter["event_id"] = options.EventID
	}

	if options.TemplateID != nil && !options.TemplateID.IsZero() {
		filter["template_id"] = options.TemplateID
	}

	var forms []*models.Form
	pagination := &PaginationOptions{
		Page:     options.Page,
		PageSize: options.PageSize,
	}

	count, err := r.mongoRepo.FindWithPagination(ctx, models.Form{}.TableName(), filter, &forms, pagination)
	if err != nil {
		return nil, 0, err
	}

	return forms, count, nil
}

// Update implements FormRepository.Update
func (r *mongoFormRepository) Update(ctx context.Context, form *models.Form) error {
	form.SetUpdatedAt(time.Now())

	filter := map[string]interface{}{
		"_id":         form.ID,
		"merchant_id": form.MerchantID,
	}

	return r.mongoRepo.UpdateOne(ctx, form.TableName(), filter, form)
}

// Delete implements FormRepository.Delete
func (r *mongoFormRepository) Delete(ctx context.Context, formID primitive.ObjectID, merchantID string) error {
	filter := map[string]interface{}{
		"_id":         formID,
		"merchant_id": merchantID,
	}

	return r.mongoRepo.DeleteOne(ctx, models.Form{}.TableName(), filter)
}

// Exists implements FormRepository.Exists
func (r *mongoFormRepository) Exists(ctx context.Context, formID primitive.ObjectID, merchantID string) (bool, error) {
	count, err := r.mongoRepo.Count(ctx, models.Form{}.TableName(), map[string]interface{}{
		"_id":         formID,
		"merchant_id": merchantID,
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByEventID implements FormRepository.FindByEventID
func (r *mongoFormRepository) FindByEventID(ctx context.Context, eventID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error) {
	filter := map[string]interface{}{
		"event_id":    eventID,
		"merchant_id": merchantID,
	}

	var forms []*models.Form
	pagination := &PaginationOptions{
		Page:     page,
		PageSize: pageSize,
	}

	count, err := r.mongoRepo.FindWithPagination(ctx, models.Form{}.TableName(), filter, &forms, pagination)
	if err != nil {
		return nil, 0, err
	}

	return forms, count, nil
}

// FindByTemplateID implements FormRepository.FindByTemplateID
func (r *mongoFormRepository) FindByTemplateID(ctx context.Context, templateID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error) {
	filter := map[string]interface{}{
		"template_id": templateID,
		"merchant_id": merchantID,
	}

	var forms []*models.Form
	pagination := &PaginationOptions{
		Page:     page,
		PageSize: pageSize,
	}

	count, err := r.mongoRepo.FindWithPagination(ctx, models.Form{}.TableName(), filter, &forms, pagination)
	if err != nil {
		return nil, 0, err
	}

	return forms, count, nil
}

// CountByTemplateID implements FormRepository.CountByTemplateID
func (r *mongoFormRepository) CountByTemplateID(ctx context.Context, templateID primitive.ObjectID, merchantID string) (int64, error) {
	filter := map[string]interface{}{
		"template_id": templateID,
		"merchant_id": merchantID,
	}

	return r.mongoRepo.Count(ctx, models.Form{}.TableName(), filter)
}