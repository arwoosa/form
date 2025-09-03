package cache

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCacheNotConnected = errors.New("cache not connected")
	ErrCacheQueryFailed  = errors.New("cache query failed")

	StatusCacheNotConnected = status.New(codes.Aborted, "cache not connected")
	StatusCacheQueryFailed  = status.New(codes.Internal, "cache query failed")
)

func ToStatus(err error) *status.Status {
	if err == nil {
		return nil
	}
	var baseSt *status.Status

	switch {
	case errors.Is(err, ErrCacheNotConnected):
		baseSt = StatusCacheNotConnected
	case errors.Is(err, ErrCacheQueryFailed):
		baseSt = StatusCacheQueryFailed
	default:
		// For unhandled errors, create a generic internal error status.
		return status.New(codes.Internal, err.Error())
	}
	unwrapErr := errors.Unwrap(err)
	if unwrapErr == nil {
		unwrapErr = err
	}
	// Add more details to the status, such as the type of violation and a description.
	st, myErr := baseSt.WithDetails(
		&errdetails.PreconditionFailure{
			Violations: []*errdetails.PreconditionFailure_Violation{
				{
					Type:        "CACHE",
					Subject:     unwrapErr.Error(),
					Description: err.Error(),
				},
			},
		},
	)
	if myErr != nil {
		// If adding details fails, return the original base status.
		return baseSt
	}
	return st
}
