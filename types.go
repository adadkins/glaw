package glaw

import (
	"net/http"

	"go.uber.org/zap"
)

type LemmyClient struct {
	baseURL   string
	APIToken  string
	jwtCookie string
	logger    *zap.Logger
	client    *http.Client
}

type Client interface {
	StreamNewPosts(pauseAfter int, closeChan chan struct{}) chan Post
	SendPrivateMessage(messageBody string, targetUserID int) error
	StreamNewComments(pauseAfter int, closeChan chan struct{}) chan Comment
	GetPost(id int) (Post, error)
	GetComment(id int) (Comment, error)
}

type PostsResponse struct {
	PostView []PostView `json:"posts"`
}

type Comment struct {
	ID        int    `json:"id"`
	CreatorID int    `json:"creator_id"`
	PostID    int    `json:"post_id"`
	Content   string `json:"content"`
	Removed   bool   `json:"removed"`
	// Published     time.Time `json:"published"`
	Deleted       bool   `json:"deleted"`
	ApID          string `json:"ap_id"`
	Local         bool   `json:"local"`
	Path          string `json:"path"`
	Distinguished bool   `json:"distinguished"`
	LanguageID    int    `json:"language_id"`
}

type CommentsResponse struct {
	Comments []Comments `json:"comments"`
}

type Post struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Body              string `json:"body"`
	CreatorID         int    `json:"creator_id"`
	CommunityID       int    `json:"community_id"`
	Removed           bool   `json:"removed"`
	Locked            bool   `json:"locked"`
	Published         string `json:"published"`
	Deleted           bool   `json:"deleted"`
	Nsfw              bool   `json:"nsfw"`
	ApID              string `json:"ap_id"`
	Local             bool   `json:"local"`
	LanguageID        int    `json:"language_id"`
	FeaturedCommunity bool   `json:"featured_community"`
	FeaturedLocal     bool   `json:"featured_local"`
}

type Comments struct {
	Comment `json:"comment"`
	Creator `json:"creator"`
}

type Creator struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Banned     bool   `json:"banned"`
	Published  string `json:"published"`
	ActorID    string `json:"actor_id"`
	Bio        string `json:"bio"`
	Local      bool   `json:"local"`
	Banner     string `json:"banner"`
	Deleted    bool   `json:"deleted"`
	Admin      bool   `json:"admin"`
	BotAccount bool   `json:"bot_account"`
	InstanceID int    `json:"instance_id"`
}

type PostResponse struct {
	PostView `json:"post_view"`
}
type PostView struct {
	Post `json:"post"`
	// Creator `json:"creator"`
	// Counts  `json:"counts"`
}

type Counts struct {
	PostID    int `json:"post_id"`
	Comments  int `json:"comments"`
	Score     int `json:"score"`
	Upvotes   int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
	// Published time.Time `json:"published"`
}

type CommentResponse struct {
	CommentView `json:"comment_view"`
}
type CommentView struct {
	Comment `json:"comment"`
	Creator `json:"creator"`
	Post    `json:"post"`
	Counts  `json:"counts"`
}
