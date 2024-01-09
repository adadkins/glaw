package glaw

import (
	"encoding/json"
	"fmt"
)

func (lc *LemmyClient) GetPost(id int) (Post, error) {
	var postResponse PostResponse
	postsBody, err := lc.callLemmyAPI("GET", fmt.Sprintf("%s%v", "post?id=", id), nil)
	if err != nil {
		return Post{}, err
	}
	err = json.Unmarshal(postsBody, &postResponse)
	if err != nil {
		return Post{}, err
	}

	return postResponse.PostView.Post, nil
}
