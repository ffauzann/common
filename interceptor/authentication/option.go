package authentication

type Option func(o *options)

type options struct {
	excludedMethods []string
}

func WithExcludedMethods(methods ...string) Option {
	return func(o *options) {
		o.excludedMethods = methods
	}
}

func evaluateOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
