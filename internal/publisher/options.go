package publisher

import (
	"net/http"
)

type option func(opts *options)

type options struct {
	httpClient *http.Client
}

func newOptions(option ...option) options {
	opts := options{}

	for _, fn := range option {
		fn(&opts)
	}
	return opts
}

func WithHttpClient(client *http.Client) option {
	return func(opts *options) {
		opts.httpClient = client
	}
}
