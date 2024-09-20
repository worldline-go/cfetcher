package consul

import "net/http"

type option struct {
	HttpClient *http.Client
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

type Option func(*option)

// WithClient sets the transport to use for the consul client
func WithClient(client *http.Client) Option {
	return func(o *option) {
		o.HttpClient = client
	}
}
