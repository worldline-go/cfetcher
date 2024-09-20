package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/worldline-go/auth"
	"github.com/worldline-go/klient"
	"golang.org/x/oauth2"
)

type Provider struct {
	auth.Provider `cfg:",squash"`

	BasicAuth BasicTransport `cfg:"basic"`
}

// NewHTTPClient returns a new http transport based on the auth service.
func NewHTTPClient(ctx context.Context, authService Provider) (*http.Client, error) {
	if authService.Active == "" {
		authService.Active = "noop"
	}

	var roundtripper func(context.Context, http.RoundTripper) (http.RoundTripper, error)

	if authService.Active == "basic" {
		slog.Info("basic auth")
		roundtripper = authService.BasicAuth.BasicRoundTripper
	} else {
		activeProvider := authService.ActiveProvider()
		if activeProvider == nil {
			return nil, fmt.Errorf("active provide not found")
		}

		clientInternal, err := klient.New(
			klient.WithDisableRetry(true),
			klient.WithDisableBaseURLCheck(true),
		)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, oauth2.HTTPClient, clientInternal.HTTP)

		oauth2Transport, err := activeProvider.NewOauth2Shared(ctx)
		if err != nil {
			return nil, fmt.Errorf("new oauth2 shared: %w", err)
		}

		roundtripper = oauth2Transport.RoundTripper
	}

	c, err := klient.New(
		klient.WithRoundTripper(roundtripper),
		klient.WithDisableRetry(true),
		klient.WithDisableBaseURLCheck(true),
	)
	if err != nil {
		return nil, err
	}

	return c.HTTP, nil
}
