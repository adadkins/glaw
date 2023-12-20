package glaw

import (
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func NewLemmyClient(url, apiToken, cookie string, client *http.Client, logger *zap.Logger) (*LemmyClient, error) {
	if url == "" {
		return nil, errors.New("url required")
	}

	if logger == nil {
		logger = zap.NewExample()
	}

	lc := LemmyClient{
		baseURL:   url,
		APIToken:  apiToken,
		jwtCookie: cookie,
		client:    client,
		logger:    logger,
	}
	return &lc, nil
}
