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

// Mock FormRepository
type MockFormRepository struct {
	mock.Mock
}

func (m *MockFormRepository) Create(ctx context.Context, form *models.Form) error {
	args := m.Called(ctx, form)
	return args.Error(0)
}

func (m *MockFormRepository) FindByID(ctx context.Context, formID primitive.ObjectID) (*models.Form, error) {
	args := m.Called(ctx, formID)
	return args.Get(0).(*models.Form), args.Error(1)
}

func (m *MockFormRepository) Find(ctx context.Context, options *models.FormQueryOptions) ([]*models.Form, int64, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]*models.Form), args.Get(1).(int64), args.Error(2)
}

func (m *MockFormRepository) Update(ctx context.Context, form *models.Form) error {
	args := m.Called(ctx, form)
	return args.Error(0)
}

func (m *MockFormRepository) Delete(ctx context.Context, formID primitive.ObjectID) error {
	args := m.Called(ctx, formID)
	return args.Error(0)
}

func (m *MockFormRepository) Exists(ctx context.Context, formID primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, formID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFormRepository) FindByEventID(ctx context.Context, eventID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error) {
	args := m.Called(ctx, eventID, merchantID, page, pageSize)
	return args.Get(0).([]*models.Form), args.Get(1).(int64), args.Error(2)
}

func (m *MockFormRepository) FindByTemplateID(ctx context.Context, templateID primitive.ObjectID, merchantID string, page, pageSize int) ([]*models.Form, int64, error) {
	args := m.Called(ctx, templateID, merchantID, page, pageSize)
	return args.Get(0).([]*models.Form), args.Get(1).(int64), args.Error(2)
}

func (m *MockFormRepository) CountByTemplateID(ctx context.Context, templateID primitive.ObjectID, merchantID string) (int64, error) {
	args := m.Called(ctx, templateID, merchantID)
	return args.Get(0).(int64), args.Error(1)
}

// Mock FormTemplateRepository
type MockFormTemplateRepository struct {
	mock.Mock
}

func (m *MockFormTemplateRepository) Create(ctx context.Context, template *models.FormTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockFormTemplateRepository) FindByID(ctx context.Context, templateID primitive.ObjectID) (*models.FormTemplate, error) {
	args := m.Called(ctx, templateID)
	return args.Get(0).(*models.FormTemplate), args.Error(1)
}

func (m *MockFormTemplateRepository) FindByMerchantID(ctx context.Context, options *models.FormTemplateQueryOptions) ([]*models.FormTemplate, int64, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]*models.FormTemplate), args.Get(1).(int64), args.Error(2)
}

func (m *MockFormTemplateRepository) Update(ctx context.Context, template *models.FormTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockFormTemplateRepository) Delete(ctx context.Context, templateID primitive.ObjectID) error {
	args := m.Called(ctx, templateID)
	return args.Error(0)
}

func (m *MockFormTemplateRepository) CountByMerchantID(ctx context.Context, merchantID string) (int64, error) {
	args := m.Called(ctx, merchantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockFormTemplateRepository) Exists(ctx context.Context, templateID primitive.ObjectID) (bool, error) {
	args := m.Called(ctx, templateID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFormTemplateRepository) Duplicate(ctx context.Context, sourceID primitive.ObjectID, nameSuffix, createdBy, merchantID string) (*models.FormTemplate, error) {
	args := m.Called(ctx, sourceID, nameSuffix, createdBy, merchantID)
	return args.Get(0).(*models.FormTemplate), args.Error(1)
}

// Test setup helper
func setupFormService() (*FormService, *MockFormRepository, *MockFormTemplateRepository, *conf.AppConfig) {
	mockFormRepo := &MockFormRepository{}
	mockTemplateRepo := &MockFormTemplateRepository{}
	config := &conf.AppConfig{
		PaginationConfig: &conf.PaginationConfig{
			DefaultPageSize: 20,
			MaxPageSize:     100,
		},
	}
	service := NewFormService(mockFormRepo, mockTemplateRepo, config)
	return service, mockFormRepo, mockTemplateRepo, config
}

// Test data helpers
func createTestForm() *models.Form {
	eventID := primitive.NewObjectID()
	return &models.Form{
		ID:         primitive.NewObjectID(),
		EventID:    &eventID,
		MerchantID: "merchant123",
		Schema:     map[string]interface{}{"type": "object"},
		UISchema:   map[string]interface{}{"ui:order": []string{}},
		CreatedBy:  "user123",
		UpdatedBy:  "user123",
	}
}

func createTestCreateFormInput() *models.CreateFormInput {
	eventID := primitive.NewObjectID()
	return &models.CreateFormInput{
		EventID:    &eventID,
		MerchantID: "merchant123",
		Schema:     map[string]interface{}{"type": "object"},
		UISchema:   map[string]interface{}{"ui:order": []string{}},
		CreatedBy:  "user123",
	}
}

func createTestUpdateFormInput() *models.UpdateFormInput {
	return &models.UpdateFormInput{
		ID:        primitive.NewObjectID(),
		Schema:    map[string]interface{}{"type": "object", "updated": true},
		UISchema:  map[string]interface{}{"ui:order": []string{"field1"}},
		UpdatedBy: "user456",
	}
}

func TestFormService_CreateForm_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	input := createTestCreateFormInput()

	mockFormRepo.On("Create", ctx, mock.AnythingOfType("*models.Form")).Return(nil)

	form, err := service.CreateForm(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, form)
	assert.Equal(t, input.EventID, form.EventID)
	assert.Equal(t, input.MerchantID, form.MerchantID)
	assert.Equal(t, input.Schema, form.Schema)
	assert.Equal(t, input.UISchema, form.UISchema)
	assert.Equal(t, input.CreatedBy, form.CreatedBy)
	assert.Equal(t, input.CreatedBy, form.UpdatedBy)
	assert.False(t, form.ID.IsZero())

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_CreateForm_ValidationError(t *testing.T) {
	service, _, _, _ := setupFormService()
	ctx := context.Background()

	invalidInput := &models.CreateFormInput{
		MerchantID: "merchant123",
		Schema:     map[string]interface{}{"type": "object"},
		UISchema:   map[string]interface{}{"ui:order": []string{}},
	}

	form, err := service.CreateForm(ctx, invalidInput)

	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestFormService_CreateForm_RepositoryError(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	input := createTestCreateFormInput()

	mockFormRepo.On("Create", ctx, mock.AnythingOfType("*models.Form")).Return(errors.New("database error"))

	form, err := service.CreateForm(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Equal(t, ErrInternalError, err)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_GetForm_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	formID := primitive.NewObjectID()
	expectedForm := createTestForm()
	expectedForm.ID = formID

	mockFormRepo.On("FindByID", ctx, formID).Return(expectedForm, nil)

	form, err := service.GetForm(ctx, formID)

	assert.NoError(t, err)
	assert.NotNil(t, form)
	assert.Equal(t, formID, form.ID)
	assert.Equal(t, expectedForm.MerchantID, form.MerchantID)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_GetForm_NotFound(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	formID := primitive.NewObjectID()

	mockFormRepo.On("FindByID", ctx, formID).Return((*models.Form)(nil), errors.New("not found"))

	form, err := service.GetForm(ctx, formID)

	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Equal(t, ErrFormNotFound, err)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListForms_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	merchantID := "merchant123"

	options := &models.FormQueryOptions{
		MerchantID: merchantID,
		Page:       1,
		PageSize:   10,
	}

	expectedForms := []*models.Form{createTestForm()}
	expectedCount := int64(1)

	mockFormRepo.On("Find", ctx, options).Return(expectedForms, expectedCount, nil)

	forms, count, err := service.ListForms(ctx, options)

	assert.NoError(t, err)
	assert.Equal(t, expectedForms, forms)
	assert.Equal(t, expectedCount, count)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListForms_WithDefaults(t *testing.T) {
	service, mockFormRepo, _, config := setupFormService()
	ctx := context.Background()
	merchantID := "merchant123"

	options := &models.FormQueryOptions{
		MerchantID: merchantID,
	}

	expectedForms := []*models.Form{}
	expectedCount := int64(0)

	mockFormRepo.On("Find", ctx, mock.MatchedBy(func(opts *models.FormQueryOptions) bool {
		return opts.Page == 1 && opts.PageSize == config.PaginationConfig.DefaultPageSize
	})).Return(expectedForms, expectedCount, nil)

	forms, count, err := service.ListForms(ctx, options)

	assert.NoError(t, err)
	assert.Equal(t, expectedForms, forms)
	assert.Equal(t, expectedCount, count)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListForms_MaxPageSize(t *testing.T) {
	service, mockFormRepo, _, config := setupFormService()
	ctx := context.Background()
	merchantID := "merchant123"

	options := &models.FormQueryOptions{
		MerchantID: merchantID,
		PageSize:   200, // exceeds max
	}

	expectedForms := []*models.Form{}
	expectedCount := int64(0)

	mockFormRepo.On("Find", ctx, mock.MatchedBy(func(opts *models.FormQueryOptions) bool {
		return opts.PageSize == config.PaginationConfig.MaxPageSize
	})).Return(expectedForms, expectedCount, nil)

	forms, count, err := service.ListForms(ctx, options)

	assert.NoError(t, err)
	assert.Equal(t, expectedForms, forms)
	assert.Equal(t, expectedCount, count)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListForms_RepositoryError(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()

	options := &models.FormQueryOptions{
		MerchantID: "merchant123",
		Page:       1,
		PageSize:   10,
	}

	mockFormRepo.On("Find", ctx, options).Return(([]*models.Form)(nil), int64(0), errors.New("database error"))

	forms, count, err := service.ListForms(ctx, options)

	assert.Error(t, err)
	assert.Nil(t, forms)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, ErrInternalError, err)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_UpdateForm_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	input := createTestUpdateFormInput()
	existingForm := createTestForm()
	existingForm.ID = input.ID

	mockFormRepo.On("FindByID", ctx, input.ID).Return(existingForm, nil)
	mockFormRepo.On("Update", ctx, mock.MatchedBy(func(form *models.Form) bool {
		schema, ok := form.Schema.(map[string]interface{})
		return ok && form.ID == input.ID &&
			form.UpdatedBy == input.UpdatedBy &&
			len(schema) == 2 // original + updated field
	})).Return(nil)

	form, err := service.UpdateForm(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, form)
	assert.Equal(t, input.ID, form.ID)
	assert.Equal(t, input.Schema, form.Schema)
	assert.Equal(t, input.UISchema, form.UISchema)
	assert.Equal(t, input.UpdatedBy, form.UpdatedBy)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_UpdateForm_ValidationError(t *testing.T) {
	service, _, _, _ := setupFormService()
	ctx := context.Background()

	invalidInput := &models.UpdateFormInput{
		Schema:   map[string]interface{}{"type": "object"},
		UISchema: map[string]interface{}{"ui:order": []string{}},
	}

	form, err := service.UpdateForm(ctx, invalidInput)

	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Contains(t, err.Error(), "invalid input")
}

func TestFormService_UpdateForm_FormNotFound(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	input := createTestUpdateFormInput()

	mockFormRepo.On("FindByID", ctx, input.ID).Return((*models.Form)(nil), errors.New("not found"))

	form, err := service.UpdateForm(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Equal(t, ErrFormNotFound, err)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_UpdateForm_RepositoryError(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	input := createTestUpdateFormInput()
	existingForm := createTestForm()
	existingForm.ID = input.ID

	mockFormRepo.On("FindByID", ctx, input.ID).Return(existingForm, nil)
	mockFormRepo.On("Update", ctx, mock.AnythingOfType("*models.Form")).Return(errors.New("database error"))

	form, err := service.UpdateForm(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Equal(t, ErrInternalError, err)

	mockFormRepo.AssertExpectations(t)
}

func TestFormService_DeleteForm_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	formID := primitive.NewObjectID()

	mockFormRepo.On("Exists", ctx, formID).Return(true, nil)
	mockFormRepo.On("Delete", ctx, formID).Return(nil)

	err := service.DeleteForm(ctx, formID)

	assert.NoError(t, err)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_DeleteForm_NotFound(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	formID := primitive.NewObjectID()

	mockFormRepo.On("Exists", ctx, formID).Return(false, nil)

	err := service.DeleteForm(ctx, formID)

	assert.Error(t, err)
	assert.Equal(t, ErrFormNotFound, err)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_DeleteForm_ExistsError(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	formID := primitive.NewObjectID()

	mockFormRepo.On("Exists", ctx, formID).Return(false, errors.New("database error"))

	err := service.DeleteForm(ctx, formID)

	assert.Error(t, err)
	assert.Equal(t, ErrInternalError, err)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_DeleteForm_DeleteError(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	formID := primitive.NewObjectID()

	mockFormRepo.On("Exists", ctx, formID).Return(true, nil)
	mockFormRepo.On("Delete", ctx, formID).Return(errors.New("database error"))

	err := service.DeleteForm(ctx, formID)

	assert.Error(t, err)
	assert.Equal(t, ErrInternalError, err)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListFormsByEvent_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	eventID := primitive.NewObjectID()
	merchantID := "merchant123"
	page := 1
	pageSize := 10

	expectedForms := []*models.Form{createTestForm()}
	expectedCount := int64(1)

	mockFormRepo.On("FindByEventID", ctx, eventID, merchantID, page, pageSize).Return(expectedForms, expectedCount, nil)

	forms, count, err := service.ListFormsByEvent(ctx, eventID, merchantID, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expectedForms, forms)
	assert.Equal(t, expectedCount, count)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListFormsByEvent_WithDefaults(t *testing.T) {
	service, mockFormRepo, _, config := setupFormService()
	ctx := context.Background()
	eventID := primitive.NewObjectID()
	merchantID := "merchant123"

	expectedForms := []*models.Form{}
	expectedCount := int64(0)

	mockFormRepo.On("FindByEventID", ctx, eventID, merchantID, 1, config.PaginationConfig.DefaultPageSize).Return(expectedForms, expectedCount, nil)

	forms, count, err := service.ListFormsByEvent(ctx, eventID, merchantID, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedForms, forms)
	assert.Equal(t, expectedCount, count)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListFormsByTemplate_Success(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	templateID := primitive.NewObjectID()
	merchantID := "merchant123"
	page := 1
	pageSize := 10

	expectedForms := []*models.Form{createTestForm()}
	expectedCount := int64(1)

	mockFormRepo.On("FindByTemplateID", ctx, templateID, merchantID, page, pageSize).Return(expectedForms, expectedCount, nil)

	forms, count, err := service.ListFormsByTemplate(ctx, templateID, merchantID, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expectedForms, forms)
	assert.Equal(t, expectedCount, count)
	mockFormRepo.AssertExpectations(t)
}

func TestFormService_ListFormsByTemplate_RepositoryError(t *testing.T) {
	service, mockFormRepo, _, _ := setupFormService()
	ctx := context.Background()
	templateID := primitive.NewObjectID()
	merchantID := "merchant123"

	mockFormRepo.On("FindByTemplateID", ctx, templateID, merchantID, 1, 20).Return(([]*models.Form)(nil), int64(0), errors.New("database error"))

	forms, count, err := service.ListFormsByTemplate(ctx, templateID, merchantID, 0, 0)

	assert.Error(t, err)
	assert.Nil(t, forms)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, ErrInternalError, err)
	mockFormRepo.AssertExpectations(t)
}
