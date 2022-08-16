package grpcclient

import (
	"context"
	util_log "github.com/cortexproject/cortex/pkg/util/log"
	cortexmiddleware "github.com/cortexproject/cortex/pkg/util/middleware"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weaveworks/common/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Instrument(requestDuration *prometheus.HistogramVec) ([]grpc.UnaryClientInterceptor, []grpc.StreamClientInterceptor) {
	return []grpc.UnaryClientInterceptor{
			httpHeaderForwardingClientInterceptor(),
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
			middleware.ClientUserHeaderInterceptor,
			cortexmiddleware.PrometheusGRPCUnaryInstrumentation(requestDuration),
		}, []grpc.StreamClientInterceptor{
			otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer()),
			middleware.StreamClientUserHeaderInterceptor,
			cortexmiddleware.PrometheusGRPCStreamInstrumentation(requestDuration),
		}
}
func httpHeaderForwardingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		headerContentsMap, ok := ctx.Value(util_log.HeaderMapContextKey).(map[string]string)
		if ok {
			for header, contents := range headerContentsMap {
				ctx = metadata.AppendToOutgoingContext(ctx, "httpheaderforwardingnames", header)
				ctx = metadata.AppendToOutgoingContext(ctx, "httpheaderforwardingcontents", contents)
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
