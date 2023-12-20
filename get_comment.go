package glaw

import (
	"encoding/json"
	"fmt"
)

func (lc *LemmyClient) GetComment(id int) (Comment, error) {
	var commentResponse CommentResponse
	postsBody, err := lc.callLemmyAPI("GET", fmt.Sprintf("%s%v", "comment?id=", id), nil)
	if err != nil {
		return Comment{}, err
	}
	err = json.Unmarshal(postsBody, &commentResponse)
	if err != nil {
		return Comment{}, err
	}

	return commentResponse.Comment, nil
}
