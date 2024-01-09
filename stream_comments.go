package glaw

import (
	"encoding/json"
	"time"
)

// comment/list?sort=New
func (lc *LemmyClient) StreamNewComments(pauseAfter int, closeChan chan struct{}) chan Comment {
	// Initialize a set to track seen items
	seenItems := make(map[int]bool)
	commentsChan := make(chan Comment, 1000)

	go func() {
		// Initialize variables for exponential backoff
		backoff := 1 * time.Second
		maxBackoff := 16 * time.Second
		backoffReset := false
		responsesWithoutNew := 0

		for {
			commentsBody, _ := lc.callLemmyAPI("GET", "comment/list?sort=New", nil)

			var postResponse CommentsResponse
			_ = json.Unmarshal(commentsBody, &postResponse)

			for _, comment := range postResponse.Comments {
				if !seenItems[comment.Comment.ID] {
					select {
					case commentsChan <- comment.Comment:
						seenItems[comment.Comment.ID] = true
						backoffReset = true
					default:
						return
					}
				}
			}

			// Pause mechanism
			if pauseAfter > 0 && backoffReset {
				responsesWithoutNew++
				if responsesWithoutNew > pauseAfter {
					// Reset backoff and responses count
					backoff = 1 * time.Second
					backoffReset = false
					responsesWithoutNew = 0
				}
			}

			// Exponential backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			// Wait for the posts channel to be closed or a timeout
			select {
			case <-closeChan:
				close(commentsChan)
			case <-time.After(backoff):
			}
		}
	}()

	return commentsChan
}
