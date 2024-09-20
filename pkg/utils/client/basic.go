package client

import (
	"context"
	"log/slog"
	"net/http"
)

// Oauth2Transport wraps oauth2.Transport to suspend CancelRequest.
type BasicTransport struct {
	transport http.RoundTripper

	Username string `cfg:"username"`
	Password string `cfg:"password"`
}

func (t *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)

	slog.Info("basic auth", slog.String("username", t.Username))

	return t.transport.RoundTrip(req)
}

// BasicRoundTripper returns a new RoundTripper that adds a basic auth header to the request.
//
// If Source is nil, returns transport as-is.
func (t BasicTransport) BasicRoundTripper(_ context.Context, transport http.RoundTripper) (http.RoundTripper, error) {
	t.transport = transport

	return &t, nil
}
