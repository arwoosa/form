package testutils

import (
	"github.com/stretchr/testify/mock"

	"github.com/arwoosa/form/internal/models"
)

// Custom matchers for testify mock

// MatchAnyFormTemplate matches any FormTemplate pointer
func MatchAnyFormTemplate() interface{} {
	return mock.MatchedBy(func(t *models.FormTemplate) bool {
		return t != nil
	})
}

// MatchAnyForm matches any Form pointer
func MatchAnyForm() interface{} {
	return mock.MatchedBy(func(f *models.Form) bool {
		return f != nil
	})
}

// MatchAnyFormTemplateSlice matches any slice of FormTemplate pointers
func MatchAnyFormTemplateSlice() interface{} {
	return mock.MatchedBy(func(templates []*models.FormTemplate) bool {
		return templates != nil
	})
}

// MatchAnyFormSlice matches any slice of Form pointers
func MatchAnyFormSlice() interface{} {
	return mock.MatchedBy(func(forms []*models.Form) bool {
		return forms != nil
	})
}

// MatchAnyObjectIDString matches any valid ObjectID string
func MatchAnyObjectIDString() interface{} {
	return mock.MatchedBy(func(id string) bool {
		return len(id) == 24 // MongoDB ObjectID hex string length
	})
}

// Generic request matchers (avoiding service package dependency)

// MatchAnyRequestWithName matches any request struct with a Name field
func MatchAnyRequestWithName() interface{} {
	return mock.MatchedBy(func(req interface{}) bool {
		return req != nil
	})
}

// MatchAnyRequestWithID matches any request struct with an ID field
func MatchAnyRequestWithID() interface{} {
	return mock.MatchedBy(func(req interface{}) bool {
		return req != nil
	})
}

// MatchFormTemplateWithName matches a FormTemplate with specific name
func MatchFormTemplateWithName(name string) interface{} {
	return mock.MatchedBy(func(t *models.FormTemplate) bool {
		return t != nil && t.Name == name
	})
}

// MatchFormTemplateWithMerchant matches a FormTemplate with specific merchant ID
func MatchFormTemplateWithMerchant(merchantID string) interface{} {
	return mock.MatchedBy(func(t *models.FormTemplate) bool {
		return t != nil && t.MerchantID == merchantID
	})
}

// MatchFormWithEventID matches a Form with specific event ID
func MatchFormWithEventID(eventID string) interface{} {
	return mock.MatchedBy(func(f *models.Form) bool {
		return f != nil && f.EventID != nil && f.EventID.Hex() == eventID
	})
}

// MatchStringSlice matches any string slice
func MatchStringSlice() interface{} {
	return mock.MatchedBy(func(slice []string) bool {
		return slice != nil
	})
}

// MatchEmptyStringSlice matches empty string slice
func MatchEmptyStringSlice() interface{} {
	return mock.MatchedBy(func(slice []string) bool {
		return slice != nil && len(slice) == 0
	})
}

// MatchNonEmptyStringSlice matches non-empty string slice
func MatchNonEmptyStringSlice() interface{} {
	return mock.MatchedBy(func(slice []string) bool {
		return len(slice) > 0
	})
}

// MatchAnyContext matches any context
func MatchAnyContext() interface{} {
	return mock.AnythingOfType("*context.emptyCtx")
}

// Helper functions for creating specific matchers

// CreateNameMatcher creates a matcher for form templates with specific name
func CreateNameMatcher(name string) func(*models.FormTemplate) bool {
	return func(t *models.FormTemplate) bool {
		return t != nil && t.Name == name
	}
}

// CreateMerchantMatcher creates a matcher for forms/templates with specific merchant ID
func CreateMerchantMatcher(merchantID string) func(*models.FormTemplate) bool {
	return func(t *models.FormTemplate) bool {
		return t != nil && t.MerchantID == merchantID
	}
}

// CreateTemplateCountMatcher creates a matcher for template slices with specific count
func CreateTemplateCountMatcher(count int) func([]*models.FormTemplate) bool {
	return func(templates []*models.FormTemplate) bool {
		return templates != nil && len(templates) == count
	}
}

// CreateFormCountMatcher creates a matcher for form slices with specific count
func CreateFormCountMatcher(count int) func([]*models.Form) bool {
	return func(forms []*models.Form) bool {
		return forms != nil && len(forms) == count
	}
}
