package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form/conf"
	"github.com/arwoosa/form/internal/models"
)

// Test setup helper for FormTemplateService
func setupFormTemplateService() (*FormTemplateService, *MockFormTemplateRepository, *conf.AppConfig) {
	mockTemplateRepo := &MockFormTemplateRepository{}
	config := &conf.AppConfig{
		PaginationConfig: &conf.PaginationConfig{
			DefaultPageSize: 20,
			MaxPageSize:     100,
		},
		BusinessRulesConfig: &conf.BusinessRulesConfig{
			MaxTemplatesPerMerchant: 10,
		},
	}
	service := NewFormTemplateService(mockTemplateRepo, config)
	return service, mockTemplateRepo, config
}

// Test data helpers for templates
func createTestFormTemplate() *models.FormTemplate {
	return &models.FormTemplate{
		ID:         primitive.NewObjectID(),
		Name:       "Test Template",
		MerchantID: "merchant123",
		Schema:     map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		UISchema:   map[string]interface{}{"ui:order": []string{}},
		CreatedBy:  "user123",
		UpdatedBy:  "user123",
	}
}

func createTestCreateFormTemplateInput() *models.CreateFormTemplateInput {
	return &models.CreateFormTemplateInput{
		Name:       "Test Template",
		MerchantID: "merchant123",
		Schema:     map[string]interface{}{"type": "object", "properties": map[string]interface{}{}},
		UISchema:   map[string]interface{}{"ui:order": []string{}},
		CreatedBy:  "user123",
	}
}

func createTestUpdateFormTemplateInput() *models.UpdateFormTemplateInput {
	return &models.UpdateFormTemplateInput{
		ID:        primitive.NewObjectID(),
		Name:      "Updated Template",
		Schema:    map[string]interface{}{"type": "object", "updated": true},
		UISchema:  map[string]interface{}{"ui:order": []string{"field1"}},
		UpdatedBy: "user456",
	}
}

func createTestDuplicateFormTemplateInput() *models.DuplicateFormTemplateInput {
	return &models.DuplicateFormTemplateInput{
		SourceID:   primitive.NewObjectID(),
		NameSuffix: " (Copy)",
		MerchantID: "merchant123",
		CreatedBy:  "user456",
	}
}

// CreateTemplate Tests
func TestFormTemplateService_CreateTemplate_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	input := createTestCreateFormTemplateInput()

	mockRepo.On("CountByMerchantID", ctx, input.MerchantID).Return(int64(5), nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.FormTemplate")).Return(nil)

	template, err := service.CreateTemplate(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, input.Name, template.Name)
	assert.Equal(t, input.MerchantID, template.MerchantID)
	assert.Equal(t, input.Schema, template.Schema)
	assert.Equal(t, input.UISchema, template.UISchema)
	assert.Equal(t, input.CreatedBy, template.CreatedBy)
	assert.False(t, template.ID.IsZero())

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_CreateTemplate_LimitExceeded(t *testing.T) {
	service, mockRepo, config := setupFormTemplateService()
	ctx := context.Background()
	input := createTestCreateFormTemplateInput()

	mockRepo.On("CountByMerchantID", ctx, input.MerchantID).Return(int64(config.BusinessRulesConfig.MaxTemplatesPerMerchant), nil)

	template, err := service.CreateTemplate(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Equal(t, ErrTemplateLimitExceeded, err)

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_CreateTemplate_ValidationError(t *testing.T) {
	service, _, _ := setupFormTemplateService()
	ctx := context.Background()

	invalidInput := &models.CreateFormTemplateInput{
		MerchantID: "merchant123",
		Schema:     map[string]interface{}{"type": "object"},
		UISchema:   map[string]interface{}{"ui:order": []string{}},
	}

	template, err := service.CreateTemplate(ctx, invalidInput)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestFormTemplateService_CreateTemplate_CountError(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	input := createTestCreateFormTemplateInput()

	mockRepo.On("CountByMerchantID", ctx, input.MerchantID).Return(int64(0), errors.New("database error"))

	template, err := service.CreateTemplate(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Equal(t, ErrInternalError, err)

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_CreateTemplate_RepositoryError(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	input := createTestCreateFormTemplateInput()

	mockRepo.On("CountByMerchantID", ctx, input.MerchantID).Return(int64(5), nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.FormTemplate")).Return(errors.New("database error"))

	template, err := service.CreateTemplate(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Equal(t, ErrInternalError, err)

	mockRepo.AssertExpectations(t)
}

// GetTemplate Tests
func TestFormTemplateService_GetTemplate_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	templateID := primitive.NewObjectID()
	expectedTemplate := createTestFormTemplate()
	expectedTemplate.ID = templateID

	mockRepo.On("FindByID", ctx, templateID).Return(expectedTemplate, nil)

	template, err := service.GetTemplate(ctx, templateID)

	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, templateID, template.ID)
	assert.Equal(t, expectedTemplate.Name, template.Name)

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_GetTemplate_NotFound(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	templateID := primitive.NewObjectID()

	mockRepo.On("FindByID", ctx, templateID).Return((*models.FormTemplate)(nil), errors.New("not found"))

	template, err := service.GetTemplate(ctx, templateID)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Equal(t, ErrTemplateNotFound, err)

	mockRepo.AssertExpectations(t)
}

// ListTemplates Tests
func TestFormTemplateService_ListTemplates_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	merchantID := "merchant123"

	options := &models.FormTemplateQueryOptions{
		MerchantID: merchantID,
		Page:       1,
		PageSize:   10,
	}

	expectedTemplates := []*models.FormTemplate{createTestFormTemplate()}
	expectedCount := int64(1)

	mockRepo.On("FindByMerchantID", ctx, options).Return(expectedTemplates, expectedCount, nil)

	templates, count, err := service.ListTemplates(ctx, options)

	assert.NoError(t, err)
	assert.Equal(t, expectedTemplates, templates)
	assert.Equal(t, expectedCount, count)

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_ListTemplates_WithDefaults(t *testing.T) {
	service, mockRepo, config := setupFormTemplateService()
	ctx := context.Background()
	merchantID := "merchant123"

	options := &models.FormTemplateQueryOptions{
		MerchantID: merchantID,
	}

	expectedTemplates := []*models.FormTemplate{}
	expectedCount := int64(0)

	mockRepo.On("FindByMerchantID", ctx, mock.MatchedBy(func(opts *models.FormTemplateQueryOptions) bool {
		return opts.Page == 1 && opts.PageSize == config.PaginationConfig.DefaultPageSize
	})).Return(expectedTemplates, expectedCount, nil)

	templates, count, err := service.ListTemplates(ctx, options)

	assert.NoError(t, err)
	assert.Equal(t, expectedTemplates, templates)
	assert.Equal(t, expectedCount, count)

	mockRepo.AssertExpectations(t)
}

// UpdateTemplate Tests
func TestFormTemplateService_UpdateTemplate_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	input := createTestUpdateFormTemplateInput()
	existingTemplate := createTestFormTemplate()
	existingTemplate.ID = input.ID

	mockRepo.On("FindByID", ctx, input.ID).Return(existingTemplate, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(template *models.FormTemplate) bool {
		return template.ID == input.ID &&
			template.Name == input.Name &&
			template.UpdatedBy == input.UpdatedBy
	})).Return(nil)

	template, err := service.UpdateTemplate(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, input.ID, template.ID)
	assert.Equal(t, input.Name, template.Name)
	assert.Equal(t, input.Schema, template.Schema)
	assert.Equal(t, input.UISchema, template.UISchema)
	assert.Equal(t, input.UpdatedBy, template.UpdatedBy)

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_UpdateTemplate_ValidationError(t *testing.T) {
	service, _, _ := setupFormTemplateService()
	ctx := context.Background()

	invalidInput := &models.UpdateFormTemplateInput{
		Name:     "Updated Template",
		Schema:   map[string]interface{}{"type": "object"},
		UISchema: map[string]interface{}{"ui:order": []string{}},
	}

	template, err := service.UpdateTemplate(ctx, invalidInput)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestFormTemplateService_UpdateTemplate_TemplateNotFound(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	input := createTestUpdateFormTemplateInput()

	mockRepo.On("FindByID", ctx, input.ID).Return((*models.FormTemplate)(nil), errors.New("not found"))

	template, err := service.UpdateTemplate(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Equal(t, ErrTemplateNotFound, err)

	mockRepo.AssertExpectations(t)
}

// DeleteTemplate Tests
func TestFormTemplateService_DeleteTemplate_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	templateID := primitive.NewObjectID()

	mockRepo.On("Exists", ctx, templateID).Return(true, nil)
	mockRepo.On("Delete", ctx, templateID).Return(nil)

	err := service.DeleteTemplate(ctx, templateID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_DeleteTemplate_NotFound(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	templateID := primitive.NewObjectID()

	mockRepo.On("Exists", ctx, templateID).Return(false, nil)

	err := service.DeleteTemplate(ctx, templateID)

	assert.Error(t, err)
	assert.Equal(t, ErrTemplateNotFound, err)
	mockRepo.AssertExpectations(t)
}

// DuplicateTemplate Tests
func TestFormTemplateService_DuplicateTemplate_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	input := createTestDuplicateFormTemplateInput()
	expectedDuplicate := createTestFormTemplate()
	expectedDuplicate.Name = "Test Template" + input.NameSuffix

	mockRepo.On("CountByMerchantID", ctx, input.MerchantID).Return(int64(5), nil)
	mockRepo.On("Duplicate", ctx, input.SourceID, input.NameSuffix, input.CreatedBy, input.MerchantID).Return(expectedDuplicate, nil)

	template, err := service.DuplicateTemplate(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, template)
	assert.Equal(t, expectedDuplicate.Name, template.Name)

	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_DuplicateTemplate_LimitExceeded(t *testing.T) {
	service, mockRepo, config := setupFormTemplateService()
	ctx := context.Background()
	input := createTestDuplicateFormTemplateInput()

	mockRepo.On("CountByMerchantID", ctx, input.MerchantID).Return(int64(config.BusinessRulesConfig.MaxTemplatesPerMerchant), nil)

	template, err := service.DuplicateTemplate(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, template)
	assert.Equal(t, ErrTemplateLimitExceeded, err)

	mockRepo.AssertExpectations(t)
}

// checkTemplateLimit Tests (internal method testing)
func TestFormTemplateService_checkTemplateLimit_Success(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	merchantID := "merchant123"

	mockRepo.On("CountByMerchantID", ctx, merchantID).Return(int64(5), nil)

	err := service.checkTemplateLimit(ctx, merchantID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_checkTemplateLimit_LimitExceeded(t *testing.T) {
	service, mockRepo, config := setupFormTemplateService()
	ctx := context.Background()
	merchantID := "merchant123"

	mockRepo.On("CountByMerchantID", ctx, merchantID).Return(int64(config.BusinessRulesConfig.MaxTemplatesPerMerchant), nil)

	err := service.checkTemplateLimit(ctx, merchantID)

	assert.Error(t, err)
	assert.Equal(t, ErrTemplateLimitExceeded, err)
	mockRepo.AssertExpectations(t)
}

func TestFormTemplateService_checkTemplateLimit_CountError(t *testing.T) {
	service, mockRepo, _ := setupFormTemplateService()
	ctx := context.Background()
	merchantID := "merchant123"

	mockRepo.On("CountByMerchantID", ctx, merchantID).Return(int64(0), errors.New("database error"))

	err := service.checkTemplateLimit(ctx, merchantID)

	assert.Error(t, err)
	assert.Equal(t, ErrInternalError, err)
	mockRepo.AssertExpectations(t)
}
