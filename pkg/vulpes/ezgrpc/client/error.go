package client

import "errors"

var (
	ErrServiceNotFound       = errors.New("service not found")
	ErrMethodNotFound        = errors.New("method not found")
	ErrInvalidRequest        = errors.New("invalid request")
	ErrConnectionFailed      = errors.New("connection failed")
	ErrFetchServerInfoFailed = errors.New("failed to fetch server info")
)
