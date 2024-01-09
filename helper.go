package glaw

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// calls an endpoint of a Lemmy instance API
func (lc *LemmyClient) callLemmyAPI(method string, endpoint string, body io.Reader) ([]byte, error) {
	startTime := time.Now()
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", lc.baseURL, endpoint), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error: %s, status code: %v, elapsedTime: %s", err.Error(), resp.StatusCode, time.Since(startTime))
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request was not ok. code: %s, body: %v, ElapsedTime: %s", resp.Status, string(respBody), time.Since(startTime))
	}

	return respBody, nil
}
