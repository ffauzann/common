package logging

import (
	"errors"

	"google.golang.org/grpc/codes"
)

type Option func(o *options)

type options struct {
	parseError bool
	mapError   map[error]codes.Code
}

// WithErrorParser parse unknown server error with unspecified internal error.
func WithErrorParser(m map[error]codes.Code) Option {
	return func(o *options) {
		o.mapError = m
		o.parseError = true
	}
}

func evaluateOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *options) getError(err error) (codes.Code, error, bool) {
	if code, ok := o.mapError[err]; ok {
		return code, err, ok
	}
	return codes.Internal, errors.New(codes.Internal.String()), false
}
