package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FormTemplate represents a reusable form template with JSON schema and UI schema
type FormTemplate struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	MerchantID  string             `bson:"merchant_id"`
	Description string             `bson:"description"`
	Schema      interface{}        `bson:"schema"`    // JSON Schema for data structure and validation
	UISchema    interface{}        `bson:"ui_schema"` // UI Schema for form layout and appearance
	CreatedAt   primitive.DateTime `bson:"created_at"`
	CreatedBy   string             `bson:"created_by"`
	UpdatedAt   primitive.DateTime `bson:"updated_at"`
	UpdatedBy   string             `bson:"updated_by"`
}

// TableName returns the collection name for FormTemplate
func (FormTemplate) TableName() string {
	return "form_templates"
}

// GetCreatedAt returns the created timestamp as time.Time
func (ft FormTemplate) GetCreatedAt() time.Time {
	return ft.CreatedAt.Time()
}

// GetUpdatedAt returns the updated timestamp as time.Time
func (ft FormTemplate) GetUpdatedAt() time.Time {
	return ft.UpdatedAt.Time()
}

// SetCreatedAt sets the created timestamp from time.Time
func (ft *FormTemplate) SetCreatedAt(t time.Time) {
	ft.CreatedAt = primitive.NewDateTimeFromTime(t)
}

// SetUpdatedAt sets the updated timestamp from time.Time
func (ft *FormTemplate) SetUpdatedAt(t time.Time) {
	ft.UpdatedAt = primitive.NewDateTimeFromTime(t)
}

// IsValid checks if the FormTemplate has required fields
func (ft FormTemplate) IsValid() bool {
	return ft.Name != "" &&
		ft.MerchantID != "" &&
		ft.CreatedBy != ""
}

// CreateFormTemplateInput represents the input for creating a new form template
type CreateFormTemplateInput struct {
	Name        string      `json:"name" validate:"required,min=1,max=100"`
	Description string      `json:"description" validate:"max=500"`
	Schema      interface{} `json:"schema" validate:"required"`
	UISchema    interface{} `json:"ui_schema"`
	CreatedBy   string      `json:"created_by" validate:"required"`
	MerchantID  string      `json:"merchant_id" validate:"required"`
}

// UpdateFormTemplateInput represents the input for updating a form template
type UpdateFormTemplateInput struct {
	ID          primitive.ObjectID `json:"id" validate:"required"`
	Name        string             `json:"name" validate:"required,min=1,max=100"`
	Description string             `json:"description" validate:"max=500"`
	Schema      interface{}        `json:"schema" validate:"required"`
	UISchema    interface{}        `json:"ui_schema"`
	UpdatedBy   string             `json:"updated_by" validate:"required"`
}

// DuplicateFormTemplateInput represents the input for duplicating a form template
type DuplicateFormTemplateInput struct {
	SourceID   primitive.ObjectID `json:"source_id" validate:"required"`
	Name       string             `json:"name" validate:"required,min=1,max=100"`
	CreatedBy  string             `json:"created_by" validate:"required"`
	MerchantID string             `json:"merchant_id" validate:"required"`
}
