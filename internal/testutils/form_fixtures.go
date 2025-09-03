package testutils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form-service/internal/models"
)

const (
	TestMerchantIDValue = "test-merchant-id"
)

// TestFormTemplate creates a test form template with default values
func TestFormTemplate() *models.FormTemplate {
	return &models.FormTemplate{
		ID:          primitive.NewObjectID(),
		Name:        "Test Form Template",
		MerchantID:  TestMerchantIDValue,
		Description: "Test form template description",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":  "string",
					"title": "Full Name",
				},
				"email": map[string]interface{}{
					"type":   "string",
					"format": "email",
					"title":  "Email Address",
				},
			},
			"required": []string{"name", "email"},
		},
		UISchema: map[string]interface{}{
			"name": map[string]interface{}{
				"ui:placeholder": "Enter your full name",
			},
			"email": map[string]interface{}{
				"ui:placeholder": "Enter your email address",
			},
		},
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		CreatedBy: "test-user-id",
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedBy: "test-user-id",
	}
}

// TestFormTemplateWithName creates a test form template with specified name
func TestFormTemplateWithName(name string) *models.FormTemplate {
	template := TestFormTemplate()
	template.Name = name
	return template
}

// TestForm creates a test form with default values
func TestForm() *models.Form {
	return &models.Form{
		ID:          primitive.NewObjectID(),
		Name:        "Test Form",
		MerchantID:  TestMerchantIDValue,
		Description: "Test form description",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"feedback": map[string]interface{}{
					"type":  "string",
					"title": "Your Feedback",
				},
				"rating": map[string]interface{}{
					"type":    "integer",
					"minimum": 1,
					"maximum": 5,
					"title":   "Rating",
				},
			},
			"required": []string{"feedback"},
		},
		UISchema: map[string]interface{}{
			"feedback": map[string]interface{}{
				"ui:widget": "textarea",
			},
			"rating": map[string]interface{}{
				"ui:widget": "updown",
			},
		},
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		CreatedBy: "test-user-id",
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedBy: "test-user-id",
	}
}

// TestFormWithTemplate creates a test form with a template ID
func TestFormWithTemplate(templateID primitive.ObjectID) *models.Form {
	form := TestForm()
	form.TemplateID = &templateID
	return form
}

// TestFormWithEventID creates a test form with an event ID
func TestFormWithEventID(eventID primitive.ObjectID) *models.Form {
	form := TestForm()
	form.EventID = &eventID
	return form
}

// TestCreateFormTemplateInput creates test input for creating a form template
func TestCreateFormTemplateInput() models.CreateFormTemplateInput {
	return models.CreateFormTemplateInput{
		Name:       "Test Template Input",
		MerchantID: TestMerchantIDValue,
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":  "string",
					"title": "Title",
				},
			},
		},
		CreatedBy: "test-user-id",
	}
}

// TestCreateFormInput creates test input for creating a form
func TestCreateFormInput() models.CreateFormInput {
	return models.CreateFormInput{
		Name:       "Test Form Input",
		MerchantID: TestMerchantIDValue,
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":  "string",
					"title": "Message",
				},
			},
		},
		CreatedBy: "test-user-id",
	}
}

// TestMerchantID returns the test merchant ID
func TestMerchantID() string {
	return TestMerchantIDValue
}

// TestUserID creates a test user ID
func TestUserID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// InvalidObjectID returns an invalid ObjectID string for testing
func InvalidObjectID() string {
	return "invalid_object_id"
}

// ValidObjectIDString returns a valid ObjectID string for testing
func ValidObjectIDString() string {
	return primitive.NewObjectID().Hex()
}
