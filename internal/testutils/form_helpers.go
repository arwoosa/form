package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/arwoosa/form/internal/models"
)

// AssertFormTemplateEqual compares two form templates for equality
func AssertFormTemplateEqual(t *testing.T, expected, actual *models.FormTemplate, ignoreTimestamps bool) {
	t.Helper()

	require.NotNil(t, expected, "Expected form template is nil")
	require.NotNil(t, actual, "Actual form template is nil")

	assert.Equal(t, expected.Name, actual.Name, "Name mismatch")
	assert.Equal(t, expected.MerchantID, actual.MerchantID, "MerchantID mismatch")
	assert.Equal(t, expected.Description, actual.Description, "Description mismatch")
	assert.Equal(t, expected.Schema, actual.Schema, "Schema mismatch")
	assert.Equal(t, expected.UISchema, actual.UISchema, "UISchema mismatch")

	if !ignoreTimestamps {
		assert.Equal(t, expected.CreatedAt, actual.CreatedAt, "CreatedAt mismatch")
		assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt, "UpdatedAt mismatch")
		assert.Equal(t, expected.CreatedBy, actual.CreatedBy, "CreatedBy mismatch")
		assert.Equal(t, expected.UpdatedBy, actual.UpdatedBy, "UpdatedBy mismatch")
	}
}

// AssertFormEqual compares two forms for equality
func AssertFormEqual(t *testing.T, expected, actual *models.Form, ignoreTimestamps bool) {
	t.Helper()

	require.NotNil(t, expected, "Expected form is nil")
	require.NotNil(t, actual, "Actual form is nil")

	assert.Equal(t, expected.Name, actual.Name, "Name mismatch")
	assert.Equal(t, expected.MerchantID, actual.MerchantID, "MerchantID mismatch")
	assert.Equal(t, expected.EventID, actual.EventID, "EventID mismatch")
	assert.Equal(t, expected.Description, actual.Description, "Description mismatch")
	assert.Equal(t, expected.Schema, actual.Schema, "Schema mismatch")
	assert.Equal(t, expected.UISchema, actual.UISchema, "UISchema mismatch")

	if !ignoreTimestamps {
		assert.Equal(t, expected.CreatedAt, actual.CreatedAt, "CreatedAt mismatch")
		assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt, "UpdatedAt mismatch")
		assert.Equal(t, expected.CreatedBy, actual.CreatedBy, "CreatedBy mismatch")
		assert.Equal(t, expected.UpdatedBy, actual.UpdatedBy, "UpdatedBy mismatch")
	}
}

// AssertValidObjectID checks that a string is a valid MongoDB ObjectID
func AssertValidObjectID(t *testing.T, id string) {
	t.Helper()

	_, err := primitive.ObjectIDFromHex(id)
	require.NoError(t, err, "Invalid ObjectID: %s", id)
}

// WaitForCondition waits for a condition to be true within a timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("Condition not met within timeout: %s", message)
}

// CreateTestContext creates a context with timeout for testing
func CreateTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// StringPtr returns a pointer to a string (helper for optional fields)
func StringPtr(s string) *string {
	return &s
}

// Int32Ptr returns a pointer to an int32 (helper for optional fields)
func Int32Ptr(i int32) *int32 {
	return &i
}

// Float64Ptr returns a pointer to a float64 (helper for optional fields)
func Float64Ptr(f float64) *float64 {
	return &f
}

// IntPtr returns a pointer to an int (helper for optional fields)
func IntPtr(i int) *int {
	return &i
}

// TimePtr returns a pointer to a time.Time (helper for optional fields)
func TimePtr(t time.Time) *time.Time {
	return &t
}
