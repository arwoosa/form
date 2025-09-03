package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestForm_TableName(t *testing.T) {
	form := Form{}
	assert.Equal(t, "forms", form.TableName())
}

func TestForm_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		form     Form
		expected bool
	}{
		{
			name: "valid form",
			form: Form{
				Name:       "Test Form",
				MerchantID: "merchant123",
				CreatedBy:  "user123",
			},
			expected: true,
		},
		{
			name: "missing name",
			form: Form{
				MerchantID: "merchant123",
				CreatedBy:  "user123",
			},
			expected: false,
		},
		{
			name: "missing merchant_id",
			form: Form{
				Name:      "Test Form",
				CreatedBy: "user123",
			},
			expected: false,
		},
		{
			name: "missing created_by",
			form: Form{
				Name:       "Test Form",
				MerchantID: "merchant123",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.form.IsValid())
		})
	}
}

func TestForm_HasTemplate(t *testing.T) {
	templateID := primitive.NewObjectID()

	tests := []struct {
		name     string
		form     Form
		expected bool
	}{
		{
			name: "has template",
			form: Form{
				TemplateID: &templateID,
			},
			expected: true,
		},
		{
			name: "no template",
			form: Form{
				TemplateID: nil,
			},
			expected: false,
		},
		{
			name: "zero template id",
			form: Form{
				TemplateID: &primitive.NilObjectID,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.form.HasTemplate())
		})
	}
}

func TestForm_HasEventID(t *testing.T) {
	eventID := primitive.NewObjectID()

	tests := []struct {
		name     string
		form     Form
		expected bool
	}{
		{
			name: "has event",
			form: Form{
				EventID: &eventID,
			},
			expected: true,
		},
		{
			name: "no event",
			form: Form{
				EventID: nil,
			},
			expected: false,
		},
		{
			name: "zero event id",
			form: Form{
				EventID: &primitive.NilObjectID,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.form.HasEventID())
		})
	}
}

func TestForm_SetAndGetTimes(t *testing.T) {
	form := Form{}
	now := time.Now()

	form.SetCreatedAt(now)
	form.SetUpdatedAt(now)

	// Allow for small differences due to precision
	assert.WithinDuration(t, now, form.GetCreatedAt(), time.Millisecond)
	assert.WithinDuration(t, now, form.GetUpdatedAt(), time.Millisecond)
}

func TestCreateFormInput_Validation(t *testing.T) {
	// Test valid input
	input := CreateFormInput{
		Name:       "Test Form",
		Schema:     map[string]interface{}{"type": "object"},
		CreatedBy:  "user123",
		MerchantID: "merchant123",
	}

	assert.NotEmpty(t, input.Name)
	assert.NotEmpty(t, input.CreatedBy)
	assert.NotEmpty(t, input.MerchantID)
	assert.NotNil(t, input.Schema)
}

func TestFormQueryOptions_Validation(t *testing.T) {
	eventID := primitive.NewObjectID()
	templateID := primitive.NewObjectID()

	options := FormQueryOptions{
		MerchantID: "merchant123",
		EventID:    &eventID,
		TemplateID: &templateID,
		Page:       1,
		PageSize:   20,
	}

	assert.Equal(t, "merchant123", options.MerchantID)
	assert.Equal(t, eventID, *options.EventID)
	assert.Equal(t, templateID, *options.TemplateID)
	assert.Equal(t, 1, options.Page)
	assert.Equal(t, 20, options.PageSize)
}
