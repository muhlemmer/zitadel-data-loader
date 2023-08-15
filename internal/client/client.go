package client

import (
	"context"
	"time"

	"github.com/muhlemmer/zitadel-data-loader/internal/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TimeoutCTX(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, config.Global.Timeout)
}

func Dial(ctx context.Context) (*grpc.ClientConn, error) {
	ctx, cancel := TimeoutCTX(ctx)
	defer cancel()

	c := config.Global

	opts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				start := time.Now()
				err := invoker(ctx, method, req, resp, cc, opts...)
				zerolog.Ctx(ctx).Err(err).TimeDiff("took_ms", time.Now(), start).Func(func(e *zerolog.Event) {
					assertProtoLog("request", req, e)
					assertProtoLog("response", resp, e)
				}).Msgf("call %q", method)
				return nil
			},
		),
	}
	if c.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	}

	cc, err := grpc.DialContext(ctx, c.Endpoint, opts...)
	zerolog.Ctx(ctx).Err(err).
		Str("endpoint", c.Endpoint).
		Bool("insecure", c.Insecure).
		Dur("timeout", c.Timeout).
		Msg("gRPC dial")
	return cc, err
}

func assertProtoLog(key string, x any, e *zerolog.Event) {
	if message, ok := x.(proto.Message); ok {
		b, err := protojson.Marshal(message)
		e.AnErr("proto_marshal_err", err)
		e.RawJSON(key, b)
		return
	}
	e.Interface(key, x)
}
