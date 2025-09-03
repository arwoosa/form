// Package codec provides a flexible framework for encoding and decoding data structures.
// It supports multiple encoding formats (GOB, MessagePack) and uses generics for type safety.
package codec

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Standardized errors for the codec package.
var (
	// ERR_UnknownCodecMethod is returned when an unsupported codec method is specified.
	ErrUnknownCodecMethod = errors.New("unknown codec method")
	// ERR_GobEncodeFailed is returned when GOB serialization fails.
	ErrGobEncodeFailed = errors.New("gob encode failed")
	// ERR_GobDecodeFailed is returned when GOB deserialization fails.
	ErrGobDecodeFailed = errors.New("gob decode failed")
	// ERR_MsgPackEncodeFailed is returned when MessagePack serialization fails.
	ErrMsgPackEncodeFailed = errors.New("msgpack encode failed")
	// ERR_MsgPackDecodeFailed is returned when MessagePack deserialization fails.
	ErrMsgPackDecodeFailed = errors.New("msgpack decode failed")
	// ERR_Base64DecodeFailed is returned when Base64 decoding of the input string fails.
	ErrBase64DecodeFailed = errors.New("base64 decode failed")

	// Status_CodecError is a pre-defined gRPC status for codec-related errors.
	Status_CodecError = status.New(codes.Internal, "codec error")
)

// ToStatus converts a codec-related wrapper error into a gRPC status.Status.
// This is useful for providing detailed, structured error information to gRPC clients.
//
// err: The wrapped error containing the original error and a descriptive message.
// Returns a gRPC status with detailed violation information, or nil if the input error is nil.
func ToStatus(err error) *status.Status {
	if err == nil {
		return nil
	}
	unwrapErr := errors.Unwrap(err)
	if unwrapErr == nil {
		unwrapErr = err
	}
	// Enhance the base codec error status with specific details from the error.
	st, myErr := Status_CodecError.WithDetails(
		&errdetails.PreconditionFailure{
			Violations: []*errdetails.PreconditionFailure_Violation{
				{
					Type:        "CODEC",
					Subject:     unwrapErr.Error(),
					Description: err.Error(),
				},
			},
		},
	)
	// If adding details fails, return the original, less specific status.
	if myErr != nil {
		return Status_CodecError
	}
	return st
}
