package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form/conf"
	"github.com/arwoosa/form/internal/dao/repository"
	"github.com/arwoosa/form/internal/models"
	"github.com/arwoosa/vulpes/log"
	"github.com/arwoosa/vulpes/relation"
	"github.com/arwoosa/vulpes/validate"
)

// FormTemplateService handles form template business logic
type FormTemplateService struct {
	templateRepo repository.FormTemplateRepository
	config       *conf.AppConfig
}

// NewFormTemplateService creates a new form template service
func NewFormTemplateService(templateRepo repository.FormTemplateRepository, config *conf.AppConfig) *FormTemplateService {
	return &FormTemplateService{
		templateRepo: templateRepo,
		config:       config,
	}
}

// CreateTemplate creates a new form template
func (s *FormTemplateService) CreateTemplate(ctx context.Context, input *models.CreateFormTemplateInput) (*models.FormTemplate, error) {
	// Validate input
	if err := validate.Struct(input); err != nil {
		log.Error("CreateTemplate validation failed", log.Err(err))
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Check template limit for merchant
	if err := s.checkTemplateLimit(ctx, input.MerchantID); err != nil {
		return nil, err
	}

	// Create template model
	template := &models.FormTemplate{
		ID:         primitive.NewObjectID(),
		Name:       input.Name,
		MerchantID: input.MerchantID,
		Schema:     input.Schema,
		UISchema:   input.UISchema,
		CreatedBy:  input.CreatedBy,
		UpdatedBy:  input.CreatedBy,
	}

	// Save to repository
	if err := s.templateRepo.Create(ctx, template); err != nil {
		log.Error("Failed to create template", log.Err(err))
		return nil, ErrInternalError
	}

	// Add Keto relation tuple for template owner
	if err := relation.AddUserResourceRole(ctx, input.CreatedBy, "FormTemplate", template.ID.Hex(), relation.RoleOwner); err != nil {
		log.Error("Failed to create Keto relation tuple for template", log.Err(err))
		// Rollback: delete the created template since Keto operation failed
		if deleteErr := s.templateRepo.Delete(ctx, template.ID); deleteErr != nil {
			log.Error("Failed to rollback template creation", log.Err(deleteErr))
		}
		return nil, fmt.Errorf("failed to create access control: %w", err)
	}

	log.Info("Template created successfully",
		log.String("template_id", template.ID.Hex()),
		log.String("name", template.Name),
		log.String("merchant_id", template.MerchantID))

	return template, nil
}

// GetTemplate retrieves a form template by ID
func (s *FormTemplateService) GetTemplate(ctx context.Context, templateID primitive.ObjectID) (*models.FormTemplate, error) {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		log.Error("Failed to get template", log.Err(err), log.String("template_id", templateID.Hex()))
		return nil, ErrTemplateNotFound
	}

	return template, nil
}

// ListTemplates retrieves form templates with pagination
func (s *FormTemplateService) ListTemplates(ctx context.Context, options *models.FormTemplateQueryOptions) ([]*models.FormTemplate, int64, error) {
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

	templates, count, err := s.templateRepo.FindByMerchantID(ctx, options)
	if err != nil {
		log.Error("Failed to list templates", log.Err(err))
		return nil, 0, ErrInternalError
	}

	return templates, count, nil
}

// UpdateTemplate updates an existing form template
func (s *FormTemplateService) UpdateTemplate(ctx context.Context, input *models.UpdateFormTemplateInput) (*models.FormTemplate, error) {
	// Validate input
	if err := validate.Struct(input); err != nil {
		log.Error("UpdateTemplate validation failed", log.Err(err))
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Get existing template to validate ownership
	existing, err := s.templateRepo.FindByID(ctx, input.ID)
	if err != nil {
		log.Error("Template not found for update", log.Err(err), log.String("template_id", input.ID.Hex()))
		return nil, ErrTemplateNotFound
	}

	// Update template fields
	existing.Name = input.Name
	existing.Schema = input.Schema
	existing.UISchema = input.UISchema
	existing.UpdatedBy = input.UpdatedBy

	// Save updates
	if err := s.templateRepo.Update(ctx, existing); err != nil {
		log.Error("Failed to update template", log.Err(err))
		return nil, ErrInternalError
	}

	log.Info("Template updated successfully",
		log.String("template_id", existing.ID.Hex()),
		log.String("name", existing.Name))

	return existing, nil
}

// DeleteTemplate deletes a form template
func (s *FormTemplateService) DeleteTemplate(ctx context.Context, templateID primitive.ObjectID) error {
	// Check if template exists
	exists, err := s.templateRepo.Exists(ctx, templateID)
	if err != nil {
		log.Error("Failed to check template existence", log.Err(err))
		return ErrInternalError
	}
	if !exists {
		return ErrTemplateNotFound
	}

	// Delete Keto relation tuples first (best effort)
	if err := relation.DeleteObjectId(ctx, "FormTemplate", templateID.Hex()); err != nil {
		log.Error("Failed to delete Keto relation tuples for template - continuing with deletion", log.Err(err))
		// Don't return here - continue with database cleanup to avoid data inconsistency
	}

	// Delete template
	if err := s.templateRepo.Delete(ctx, templateID); err != nil {
		log.Error("Failed to delete template", log.Err(err))
		return ErrInternalError
	}

	log.Info("Template deleted successfully",
		log.String("template_id", templateID.Hex()))

	return nil
}

// DuplicateTemplate creates a duplicate of an existing template
func (s *FormTemplateService) DuplicateTemplate(ctx context.Context, input *models.DuplicateFormTemplateInput) (*models.FormTemplate, error) {
	// Validate input
	if err := validate.Struct(input); err != nil {
		log.Error("DuplicateTemplate validation failed", log.Err(err))
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Check template limit for merchant
	if err := s.checkTemplateLimit(ctx, input.MerchantID); err != nil {
		return nil, err
	}

	// Duplicate template
	duplicate, err := s.templateRepo.Duplicate(ctx, input.SourceID, input.NameSuffix, input.CreatedBy, input.MerchantID)
	if err != nil {
		log.Error("Failed to duplicate template", log.Err(err))
		return nil, ErrInternalError
	}

	// Add Keto relation tuple for duplicated template owner
	if err := relation.AddUserResourceRole(ctx, input.CreatedBy, "FormTemplate", duplicate.ID.Hex(), relation.RoleOwner); err != nil {
		log.Error("Failed to create Keto relation tuple for duplicated template", log.Err(err))
		// Rollback: delete the duplicated template since Keto operation failed
		if deleteErr := s.templateRepo.Delete(ctx, duplicate.ID); deleteErr != nil {
			log.Error("Failed to rollback template duplication", log.Err(deleteErr))
		}
		return nil, fmt.Errorf("failed to create access control: %w", err)
	}

	log.Info("Template duplicated successfully",
		log.String("source_id", input.SourceID.Hex()),
		log.String("new_id", duplicate.ID.Hex()),
		log.String("new_name", duplicate.Name))

	return duplicate, nil
}

// checkTemplateLimit validates if merchant can create more templates
func (s *FormTemplateService) checkTemplateLimit(ctx context.Context, merchantID string) error {
	count, err := s.templateRepo.CountByMerchantID(ctx, merchantID)
	if err != nil {
		log.Error("Failed to count templates", log.Err(err))
		return ErrInternalError
	}

	if count >= int64(s.config.BusinessRulesConfig.MaxTemplatesPerMerchant) {
		log.Warn("Template limit exceeded",
			log.String("merchant_id", merchantID),
			log.Int64("current_count", count),
			log.Int("limit", s.config.BusinessRulesConfig.MaxTemplatesPerMerchant))
		return ErrTemplateLimitExceeded
	}

	return nil
}
