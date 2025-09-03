package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Form represents an individual form instance that can be based on a template or have custom schema
type Form struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	Name        string              `bson:"name"`
	EventID     *primitive.ObjectID `bson:"event_id,omitempty"` // Optional reference to an event
	MerchantID  string              `bson:"merchant_id"`
	TemplateID  *primitive.ObjectID `bson:"template_id,omitempty"` // Optional reference to a form template
	Description string              `bson:"description"`
	Schema      interface{}         `bson:"schema"`    // JSON Schema for data structure and validation
	UISchema    interface{}         `bson:"ui_schema"` // UI Schema for form layout and appearance
	CreatedAt   primitive.DateTime  `bson:"created_at"`
	CreatedBy   string              `bson:"created_by"`
	UpdatedAt   primitive.DateTime  `bson:"updated_at"`
	UpdatedBy   string              `bson:"updated_by"`
}

// TableName returns the collection name for Form
func (Form) TableName() string {
	return "forms"
}

// GetCreatedAt returns the created timestamp as time.Time
func (f Form) GetCreatedAt() time.Time {
	return f.CreatedAt.Time()
}

// GetUpdatedAt returns the updated timestamp as time.Time
func (f Form) GetUpdatedAt() time.Time {
	return f.UpdatedAt.Time()
}

// SetCreatedAt sets the created timestamp from time.Time
func (f *Form) SetCreatedAt(t time.Time) {
	f.CreatedAt = primitive.NewDateTimeFromTime(t)
}

// SetUpdatedAt sets the updated timestamp from time.Time
func (f *Form) SetUpdatedAt(t time.Time) {
	f.UpdatedAt = primitive.NewDateTimeFromTime(t)
}

// IsValid checks if the Form has required fields
func (f Form) IsValid() bool {
	return f.Name != "" &&
		f.MerchantID != "" &&
		f.CreatedBy != ""
}

// HasTemplate checks if the form is based on a template
func (f Form) HasTemplate() bool {
	return f.TemplateID != nil && !f.TemplateID.IsZero()
}

// HasEventID checks if the form is associated with an event
func (f Form) HasEventID() bool {
	return f.EventID != nil && !f.EventID.IsZero()
}

// CreateFormInput represents the input for creating a new form
type CreateFormInput struct {
	Name        string              `json:"name" validate:"required,min=1,max=100"`
	EventID     *primitive.ObjectID `json:"event_id,omitempty"`
	TemplateID  *primitive.ObjectID `json:"template_id,omitempty"`
	Description string              `json:"description" validate:"max=500"`
	Schema      interface{}         `json:"schema" validate:"required"`
	UISchema    interface{}         `json:"ui_schema"`
	CreatedBy   string              `json:"created_by" validate:"required"`
	MerchantID  string              `json:"merchant_id" validate:"required"`
}

// UpdateFormInput represents the input for updating a form
type UpdateFormInput struct {
	ID          primitive.ObjectID  `json:"id" validate:"required"`
	Name        string              `json:"name" validate:"required,min=1,max=100"`
	EventID     *primitive.ObjectID `json:"event_id,omitempty"`
	TemplateID  *primitive.ObjectID `json:"template_id,omitempty"`
	Description string              `json:"description" validate:"max=500"`
	Schema      interface{}         `json:"schema" validate:"required"`
	UISchema    interface{}         `json:"ui_schema"`
	UpdatedBy   string              `json:"updated_by" validate:"required"`
}

// FormQueryOptions represents query options for listing forms
type FormQueryOptions struct {
	MerchantID string              `json:"merchant_id" validate:"required"`
	EventID    *primitive.ObjectID `json:"event_id,omitempty"`
	TemplateID *primitive.ObjectID `json:"template_id,omitempty"`
	Page       int                 `json:"page" validate:"min=1"`
	PageSize   int                 `json:"page_size" validate:"min=1,max=100"`
}

// FormTemplateQueryOptions represents query options for listing form templates
type FormTemplateQueryOptions struct {
	MerchantID string `json:"merchant_id" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}
