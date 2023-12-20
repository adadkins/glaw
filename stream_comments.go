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
			commentsBody, err := lc.callLemmyAPI("GET", "comment/list?sort=New", nil)
			if err != nil {
				lc.logger.Error(err.Error())
			}
			var postResponse CommentsResponse
			err = json.Unmarshal(commentsBody, &postResponse)
			if err != nil {
				lc.logger.Sugar().Infof("postsBody: %s", commentsBody)
				lc.logger.Error(err.Error())
			}

			for _, comment := range postResponse.Comments {
				if !seenItems[comment.Comment.ID] {
					select {
					case commentsChan <- comment.Comment:
						seenItems[comment.Comment.ID] = true
						backoffReset = true
					default:
						lc.logger.Error("Channel closed or not ready")
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
				lc.logger.Info("Comments channel closed.")
				close(commentsChan)
			case <-time.After(backoff):
			}
		}
	}()

	return commentsChan
}
