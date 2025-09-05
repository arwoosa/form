package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form/conf"
	"github.com/arwoosa/form/internal/dao/repository"
	"github.com/arwoosa/form/internal/models"
	"github.com/arwoosa/vulpes/log"
	"github.com/arwoosa/vulpes/validate"
)

// FormService handles form business logic
type FormService struct {
	formRepo     repository.FormRepository
	templateRepo repository.FormTemplateRepository
	config       *conf.AppConfig
}

// NewFormService creates a new form service
func NewFormService(formRepo repository.FormRepository, templateRepo repository.FormTemplateRepository, config *conf.AppConfig) *FormService {
	return &FormService{
		formRepo:     formRepo,
		templateRepo: templateRepo,
		config:       config,
	}
}

// CreateForm creates a new form
func (s *FormService) CreateForm(ctx context.Context, input *models.CreateFormInput) (*models.Form, error) {
	// Validate input
	if err := validate.Struct(input); err != nil {
		log.Error("CreateForm validation failed", log.Err(err))
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Create form model
	form := &models.Form{
		ID:         primitive.NewObjectID(),
		EventID:    input.EventID,
		MerchantID: input.MerchantID,
		Schema:     input.Schema,
		UISchema:   input.UISchema,
		CreatedBy:  input.CreatedBy,
		UpdatedBy:  input.CreatedBy,
	}

	// Save to repository
	if err := s.formRepo.Create(ctx, form); err != nil {
		log.Error("Failed to create form", log.Err(err))
		return nil, ErrInternalError
	}

	log.Info("Form created successfully",
		log.String("form_id", form.ID.Hex()),
		log.String("merchant_id", form.MerchantID))

	return form, nil
}

// GetForm retrieves a form by ID
func (s *FormService) GetForm(ctx context.Context, formID primitive.ObjectID) (*models.Form, error) {
	form, err := s.formRepo.FindByID(ctx, formID)
	if err != nil {
		log.Error("Failed to get form", log.Err(err), log.String("form_id", formID.Hex()))
		return nil, ErrFormNotFound
	}

	return form, nil
}

// ListForms retrieves forms with pagination and optional filters
func (s *FormService) ListForms(ctx context.Context, options *models.FormQueryOptions) ([]*models.Form, int64, error) {
	// Set default pagination if not provided
	if options.Page <= 0 {
		options.Page = 1
	}
	if options.PageSize <= 0 {
		options.PageSize = s.config.PaginationConfig.DefaultPageSize
	}
	if options.PageSize > s.config.PaginationConfig.MaxPageSize {
		options.PageSize = s.config.PaginationConfig.MaxPageSize
	}

	forms, count, err := s.formRepo.Find(ctx, options)
	if err != nil {
		log.Error("Failed to list forms", log.Err(err))
		return nil, 0, ErrInternalError
	}

	return forms, count, nil
}

// UpdateForm updates an existing form
func (s *FormService) UpdateForm(ctx context.Context, input *models.UpdateFormInput) (*models.Form, error) {
	// Validate input
	if err := validate.Struct(input); err != nil {
		log.Error("UpdateForm validation failed", log.Err(err))
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Get existing form to validate ownership
	existing, err := s.formRepo.FindByID(ctx, input.ID)
	if err != nil {
		log.Error("Form not found for update", log.Err(err), log.String("form_id", input.ID.Hex()))
		return nil, ErrFormNotFound
	}

	// Update form fields
	existing.Schema = input.Schema
	existing.UISchema = input.UISchema
	existing.UpdatedBy = input.UpdatedBy

	// Save updates
	if err := s.formRepo.Update(ctx, existing); err != nil {
		log.Error("Failed to update form", log.Err(err))
		return nil, ErrInternalError
	}

	log.Info("Form updated successfully",
		log.String("form_id", existing.ID.Hex()))

	return existing, nil
}

// DeleteForm deletes a form
func (s *FormService) DeleteForm(ctx context.Context, formID primitive.ObjectID) error {
	// Check if form exists
	exists, err := s.formRepo.Exists(ctx, formID)
	if err != nil {
		log.Error("Failed to check form existence", log.Err(err))
		return ErrInternalError
	}
	if !exists {
		return ErrFormNotFound
	}

	// Delete form
	if err := s.formRepo.Delete(ctx, formID); err != nil {
		log.Error("Failed to delete form", log.Err(err))
		return ErrInternalError
	}

	log.Info("Form deleted successfully",
		log.String("form_id", formID.Hex()))

	return nil
}

// ListFormsByEvent retrieves forms associated with an event
func (s *FormService) ListFormsByEvent(ctx context.Context, eventID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error) {
	// Set default pagination
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = s.config.PaginationConfig.DefaultPageSize
	}
	if pageSize > s.config.PaginationConfig.MaxPageSize {
		pageSize = s.config.PaginationConfig.MaxPageSize
	}

	forms, count, err := s.formRepo.FindByEventID(ctx, eventID, merchantID, page, pageSize)
	if err != nil {
		log.Error("Failed to list forms by event", log.Err(err))
		return nil, 0, ErrInternalError
	}

	return forms, count, nil
}

// ListFormsByTemplate retrieves forms associated with a template
func (s *FormService) ListFormsByTemplate(ctx context.Context, templateID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error) {
	// Set default pagination
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = s.config.PaginationConfig.DefaultPageSize
	}
	if pageSize > s.config.PaginationConfig.MaxPageSize {
		pageSize = s.config.PaginationConfig.MaxPageSize
	}

	forms, count, err := s.formRepo.FindByTemplateID(ctx, templateID, merchantID, page, pageSize)
	if err != nil {
		log.Error("Failed to list forms by template", log.Err(err))
		return nil, 0, ErrInternalError
	}

	return forms, count, nil
}
