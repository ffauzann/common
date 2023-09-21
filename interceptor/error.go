package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errCodes map[error]codes.Code

type errorWrapper struct {
	Code    codes.Code
	Message string
}

// Error implements required method to fulfill error interface.
func (e *errorWrapper) Error() string {
	return e.Message
}

// GRPCStatus implements required method to fulfill anonym interface.
// Used in status.FromError().
func (e *errorWrapper) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}

// Deprecated: use logging module instead.
// ErrorHandler handles proper error codes & messages so the error details won't be thrown to grpc client.
// e.g sql errors, pointer errors, marshall errors etc.
func ErrorHandler(mapErrCodes map[error]codes.Code) grpc.UnaryServerInterceptor {
	errCodes = mapErrCodes
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			err = func(err error) error {
				code, err := getError(err)
				return &errorWrapper{
					Code:    code,
					Message: err.Error(),
				}
			}(err)
		}

		return
	}
}

func getError(err error) (codes.Code, error) {
	if code, ok := errCodes[err]; ok {
		return code, err
	}
	return codes.Internal, errors.New(codes.Internal.String())
}
