package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFormTemplate_TableName(t *testing.T) {
	template := FormTemplate{}
	assert.Equal(t, "form_templates", template.TableName())
}

func TestFormTemplate_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		template FormTemplate
		expected bool
	}{
		{
			name: "valid template",
			template: FormTemplate{
				Name:       "Test Template",
				MerchantID: "merchant123",
				CreatedBy:  "user123",
			},
			expected: true,
		},
		{
			name: "missing name",
			template: FormTemplate{
				MerchantID: "merchant123",
				CreatedBy:  "user123",
			},
			expected: false,
		},
		{
			name: "missing merchant_id",
			template: FormTemplate{
				Name:      "Test Template",
				CreatedBy: "user123",
			},
			expected: false,
		},
		{
			name: "missing created_by",
			template: FormTemplate{
				Name:       "Test Template",
				MerchantID: "merchant123",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.template.IsValid())
		})
	}
}

func TestFormTemplate_SetAndGetTimes(t *testing.T) {
	template := FormTemplate{}
	now := time.Now()

	template.SetCreatedAt(now)
	template.SetUpdatedAt(now)

	// Allow for small differences due to precision
	assert.WithinDuration(t, now, template.GetCreatedAt(), time.Millisecond)
	assert.WithinDuration(t, now, template.GetUpdatedAt(), time.Millisecond)
}

func TestCreateFormTemplateInput_Validation(t *testing.T) {
	// Test valid input
	input := CreateFormTemplateInput{
		Name:       "Test Template",
		Schema:     map[string]interface{}{"type": "object"},
		CreatedBy:  "user123",
		MerchantID: "merchant123",
	}

	// These would typically be tested with a validator, but we can at least check fields exist
	assert.NotEmpty(t, input.Name)
	assert.NotEmpty(t, input.CreatedBy)
	assert.NotEmpty(t, input.MerchantID)
	assert.NotNil(t, input.Schema)
}

func TestDuplicateFormTemplateInput_Fields(t *testing.T) {
	sourceID := primitive.NewObjectID()

	input := DuplicateFormTemplateInput{
		SourceID:   sourceID,
		NameSuffix: "Copy",
		CreatedBy:  "user123",
		MerchantID: "merchant123",
	}

	assert.Equal(t, sourceID, input.SourceID)
	assert.Equal(t, "Copy", input.NameSuffix)
	assert.Equal(t, "user123", input.CreatedBy)
	assert.Equal(t, "merchant123", input.MerchantID)
}
