package service

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// Common errors
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternalError = errors.New("internal server error")

	// Template-specific errors
	ErrTemplateNotFound      = errors.New("form template not found")
	ErrTemplateLimitExceeded = errors.New("template limit exceeded for merchant")
	ErrTemplateNameExists    = errors.New("template name already exists")

	// Form-specific errors
	ErrFormNotFound        = errors.New("form not found")
	ErrFormInvalidTemplate = errors.New("invalid form template reference")
	ErrFormInvalidEvent    = errors.New("invalid event reference")
)

// ToGRPCError converts service errors to gRPC status errors
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case ErrUnauthorized:
		return status.Error(codes.Unauthenticated, err.Error())
	case ErrNotFound, ErrTemplateNotFound, ErrFormNotFound:
		return status.Error(codes.NotFound, err.Error())
	case ErrInvalidInput, ErrFormInvalidTemplate, ErrFormInvalidEvent:
		return status.Error(codes.InvalidArgument, err.Error())
	case ErrTemplateLimitExceeded:
		return status.Error(codes.ResourceExhausted, err.Error())
	case ErrTemplateNameExists:
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

// ValidationError represents input validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// BusinessRuleError represents business logic validation errors
type BusinessRuleError struct {
	Rule    string
	Message string
}

func (e BusinessRuleError) Error() string {
	return fmt.Sprintf("business rule violation '%s': %s", e.Rule, e.Message)
}
