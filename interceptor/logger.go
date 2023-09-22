package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/ffauzann/common/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Deprecated: Use `interceptor/logger` instead.
func Logger(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		method := util.GetMethod(info.FullMethod)
		resp, err = handler(ctx, req)

		func() {
			d := time.Since(start)
			msg := fmt.Sprintf("%s took %v", method, d)

			if ms := d.Milliseconds(); ms < 10000 {
				msg = fmt.Sprintf("%s took %vms", method, ms)
			}

			if err != nil {
				code, _ := getError(err)
				l.Info(
					msg,
					zap.Any("code", code),
					zap.Any("req", util.StructToMap(req)),
					zap.Any("res", util.StructToMap(resp)),
				)
				return
			}

			l.Info(msg)
		}()

		return
	}
}
