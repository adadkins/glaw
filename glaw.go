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

	// set default timeout
	lc := LemmyClient{
		baseURL:   url,
		APIToken:  apiToken,
		jwtCookie: cookie,
		client:    client,
		logger:    logger,
		timeout:   10,
	}
	return &lc, nil
}

func (lc *LemmyClient) SetTimeout(timeout int) {
	if timeout > 0 {
		lc.timeout = timeout
	}
}
