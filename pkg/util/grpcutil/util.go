package grpcutil

import (
	"context"
	util_log "github.com/cortexproject/cortex/pkg/util/log"
	"github.com/gogo/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// IsGRPCContextCanceled returns whether the input error is a GRPC error wrapping
// the context.Canceled error.
func IsGRPCContextCanceled(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}

	return s.Code() == codes.Canceled
}

func HTTPHeaderForwardingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		headerMap := make(map[string]string)
		meta, worked := metadata.FromIncomingContext(ctx)
		if worked {
			headersSlice := meta["httpheaderforwardingnames"]
			headerContentsSlice := meta["httpheaderforwardingcontents"]
			if len(headersSlice) == len(headerContentsSlice) {
				for i, header := range headersSlice {
					headerMap[header] = headerContentsSlice[i]
				}
				ctx = context.WithValue(ctx, util_log.HeaderMapContextKey, headerMap)
			}
		}
		h, err := handler(ctx, req)
		return h, err
	}
}
