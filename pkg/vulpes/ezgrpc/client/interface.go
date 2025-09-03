package client

import (
	"context"
)

type Client interface {
	Invoke(ctx context.Context, address, serviceName, methodName string, req []byte) ([]byte, error)
	GetServiceInvoker(ctx context.Context, address, serviceName string) (ServiceInvoker, error)
	Close() error
}

type ServiceInvoker interface {
	Invoke(ctx context.Context, methodName string, req []byte) ([]byte, error)
	IsMethodExists(methodName string) bool
}
