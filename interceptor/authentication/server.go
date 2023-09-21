package authentication

import (
	"context"
	"hangoutin/common/util"
	"strings"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	const expectedScheme = "Bearer"
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if slices.Contains(o.excludedMethods, util.GetMethod(info.FullMethod)) {
			return handler(ctx, req)
		}

		md, _ := metadata.FromIncomingContext(ctx)
		val := md.Get("authorization")
		if len(val) == 0 {
			return nil, status.Error(codes.Unauthenticated, "")
		}

		s := strings.SplitN(val[0], " ", 2)
		if len(s) < 2 {
			return nil, status.Error(codes.Unauthenticated, "")
		}

		if s[0] != expectedScheme {
			return nil, status.Error(codes.Unauthenticated, "")
		}
		return handler(ctx, req)
	}
}
