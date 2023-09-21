package logging

import (
	"context"
	"hangoutin/common/util"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(l *zap.Logger, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		method := util.GetMethod(info.FullMethod)
		resp, err = handler(ctx, req)
		msg := formatMessage(start, method)

		if err != nil {
			err = func(err error) error {
				code, maskedErr, ok := o.getError(err)
				logError(l, msg, ok, code, err)

				// Handle known error.
				if ok {
					return &errorResponse{
						Code:    code,
						Message: err.Error(),
					}
				}

				// Handle unknown error.
				return &errorResponse{
					Code:    code,
					Message: maskedErr.Error(),
				}
			}(err)

			return
		}

		// Only log success requests.
		l.Info(msg)

		return
	}
}
