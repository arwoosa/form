// Package codec provides a flexible framework for encoding and decoding data structures.
// It supports multiple encoding formats (GOB, MessagePack) and uses generics for type safety.
package codec

import (
	"sync"
)

// once ensures that the codec method can only be set once during the application's lifecycle.
var once sync.Once

// WithCodecMethod sets the global default encoding method for the package.
// This function uses sync.Once to ensure that the codec method can only be set once,
// preventing inconsistent encoding/decoding formats during runtime.
// It should be called during the application's initialization phase.
//
// Example:
//
//	func main() {
//	    codec.WithCodecMethod(codec.MSGPACK)
//	    // ... rest of the application
//	}
func WithCodecMethod(method CodecMethod) {
	once.Do(func() {
		defaultCodeMethod = method
	})
}
