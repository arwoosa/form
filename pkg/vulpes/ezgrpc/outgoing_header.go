package ezgrpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var OutgoingHeaderMatcher = runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcherHandler)

var invalidHeaderKey = []string{"content-type", "trailer"}

func outgoingHeaderMatcherHandler(key string) (string, bool) {
	for _, invalidKey := range invalidHeaderKey {
		if key == invalidKey {
			return key, false
		}
	}
	// 對其他所有 metadata，使用預設的行為
	return key, true
}
