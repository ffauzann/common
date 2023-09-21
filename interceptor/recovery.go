package interceptor

import (
	"context"
	"fmt"
	"hangoutin/common/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Deprecated: Use `interceptor/recovery` instead.
func Recovery(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				l.Error(fmt.Sprintf("%s PANIC: %v", util.GetMethod(info.FullMethod), r))
				err = status.Errorf(codes.Internal, "%v", r)
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false
		return resp, err
	}
}
