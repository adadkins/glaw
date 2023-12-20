package glaw

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// calls an endpoint of a Lemmy instance API
func (lc *LemmyClient) callLemmyAPI(method string, endpoint string, body io.Reader) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Prepare the request
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", lc.baseURL, endpoint), body)
	if err != nil {
		return nil, err
	}

	// Set the API token for authentication (if required)
	if lc.APIToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.APIToken))
	}
	if lc.jwtCookie != "" {
		req.Header.Add("cookie", lc.jwtCookie)
	}

	// Send the request
	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		lc.logger.Error(err.Error())
		return nil, err
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		lc.logger.Sugar().Infof("request was not ok. code: %s, body: %s", resp.Status, respBody)
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return respBody, nil
}
