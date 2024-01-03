package glaw

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// calls an endpoint of a Lemmy instance API
func (lc *LemmyClient) callLemmyAPI(method string, endpoint string, body io.Reader) ([]byte, error) {
	startTime := time.Now()
	// Set the maximum number of retries
	maxRetries := 5

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create a context with timeout for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lc.timeout)*time.Second)

		// Prepare the request
		req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", lc.baseURL, endpoint), body)
		if err != nil {
			cancel()
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
			cancel()

			// Check for timeout error and retry
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}

			return nil, err
		}

		// Check if the context has timed out
		select {
		case <-ctx.Done():
			cancel()
			lc.logRequestInfo(startTime, "Request timed out")
			return nil, ctx.Err()
		default:
			// Continue processing if the context has not timed out
		}
		defer func() {
			resp.Body.Close()
			cancel() // Ensure cancellation is called when closing the response body
		}()

		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lc.logger.Sugar().Infof("Error: %s, status code: %s", err.Error(), "status code: %v", resp.StatusCode)
			return nil, err
		}

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			lc.logRequestInfo(startTime, "Request failed:", resp.Status, string(respBody))
			lc.logger.Sugar().Infof("request was not ok. code: %s, body: %v", resp.Status, respBody)

			// Retry for certain HTTP status codes if needed
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusInternalServerError {
				continue
			}

			return nil, fmt.Errorf("request failed with status: %s", resp.Status)
		}

		return respBody, nil
	}

	// All attempts failed, return an error
	return nil, fmt.Errorf("maximum number of retries reached")
}

func (lc *LemmyClient) logRequestInfo(startTime time.Time, messages ...interface{}) {
	elapsed := time.Since(startTime)
	lc.logger.Sugar().Infow("Request Info",
		"ElapsedTime", elapsed,
		"Timestamp", startTime.Format(time.RFC3339),
		"Messages", messages,
	)
}
