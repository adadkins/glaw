package glaw

import (
	"errors"
	"net/http"
)

func NewLemmyClient(url, apiToken, cookie string, client *http.Client) (*LemmyClient, error) {
	if url == "" {
		return nil, errors.New("url required")
	}

	// set default timeout
	lc := LemmyClient{
		baseURL:   url,
		APIToken:  apiToken,
		jwtCookie: cookie,
		client:    client,
	}
	return &lc, nil
}
