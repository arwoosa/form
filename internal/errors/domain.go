package errors

import "errors"

// Common domain errors - these represent business rule violations or entity states
var (
	ErrFormNotFound      = errors.New("form not found")
	ErrTemplateNotFound  = errors.New("form template not found")
	ErrInvalidSchema     = errors.New("invalid form schema")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrInvalidMerchantID = errors.New("invalid merchant_id")
)
