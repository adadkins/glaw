package glaw_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/adadkins/glaw"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

// test instance and creditials
const (
	url = "https://voyager.lemmy.ml/api/v3/"
	jwt = "jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMyIsImlzcyI6InZveWFnZXIubGVtbXkubWwiLCJpYXQiOjE3MDAyNTI5MDZ9.fwYfeaOPz4G5qtZlfZgh-JLCVQJlGytVqvm0MOh3vvY"
)

func TestStreamPosts(t *testing.T) {
	t.Run("Posts channel streams posts", func(t *testing.T) {
		//given
		client, recorder := GetRecorderClient(t, "posts")
		defer recorder.Stop()

		lc, _ := glaw.NewLemmyClient(url, "", jwt, client)

		//when
		posts := lc.StreamNewPosts(5, nil)

		// pull posts out of the channel
		postsLists := []glaw.Post{}
		for i := 0; i < 10; i++ {
			p := <-posts
			postsLists = append(postsLists, p)

		}

		// construct our expected slice
		// TODO: move this into a testing values file? or fixtures?
		var expectedPosts []glaw.Post
		inputString := `[{"id":39509,"name":"XD why is invalid","body":"Post","creator_id":83995,"community_id":2738,"removed":false,"locked":false,"published":"2023-12-10T18:26:52.100597Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/39509","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":39204,"name":"This is a test post being edited now","body":"Hi friend we are editing this post is now nsfw","creator_id":83995,"community_id":2738,"removed":false,"locked":false,"published":"2023-12-10T10:37:46.649591Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/39204","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":30391,"name":"Unpopular opinion: this community is the worst","body":"Fight me","creator_id":1537,"community_id":121,"removed":false,"locked":false,"published":"2023-12-01T20:08:03.355973Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/30391","local":true,"language_id":37,"featured_community":false,"featured_local":false},{"id":28738,"name":"test image","body":"","creator_id":2086,"community_id":1501,"removed":false,"locked":false,"published":"2023-11-30T11:13:56.994682Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/28738","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":26629,"name":"Test image post","body":"","creator_id":2086,"community_id":1501,"removed":false,"locked":false,"published":"2023-11-28T19:08:25.244778Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/26629","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":26319,"name":"local post","body":"","creator_id":4613,"community_id":120,"removed":false,"locked":false,"published":"2023-11-28T13:28:44.196156Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/26319","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":16945,"name":"test image 2","body":"","creator_id":2,"community_id":120,"removed":false,"locked":false,"published":"2023-11-23T10:31:07.616307Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/16945","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":16943,"name":"test image","body":"","creator_id":2,"community_id":120,"removed":false,"locked":false,"published":"2023-11-23T10:28:43.040040Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/16943","local":true,"language_id":0,"featured_community":false,"featured_local":false},{"id":16254,"name":"hi!","body":"","creator_id":3735,"community_id":1000,"removed":false,"locked":false,"published":"2023-11-22T20:33:28.324375Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/16254","local":true,"language_id":0,"featured_community":true,"featured_local":false},{"id":16141,"name":"Test post 2","body":"Testing","creator_id":2086,"community_id":1501,"removed":false,"locked":false,"published":"2023-11-22T18:46:39.914714Z","deleted":false,"nsfw":false,"ap_id":"https://voyager.lemmy.ml/post/16141","local":true,"language_id":0,"featured_community":false,"featured_local":false}]`
		err := json.Unmarshal([]byte(inputString), &expectedPosts)
		if err != nil {
			t.Fatal(err)
		}

		//then
		assert.NotNil(t, postsLists)
		assert.Equal(t, expectedPosts, postsLists)
	})
	t.Run("Closing shuts down Posts channel", func(t *testing.T) {
		// given
		client, recorder := GetRecorderClient(t, "posts")
		defer recorder.Stop()
		lc, _ := glaw.NewLemmyClient(url, "", jwt, client)
		close := make(chan struct{})

		//when
		postsChan := lc.StreamNewPosts(5, close)

		time.Sleep(1 * time.Second)

		// send done signal
		close <- struct{}{}

		// pull all items out of channel
		for {
			_, ok := <-postsChan
			if !ok {
				// Channel closed, and all values are received
				break
			}
		}

		// blocking check if channel is closed
		_, ok := <-postsChan

		// then
		assert.False(t, ok, "Expected posts channel to be closed, but it's still open.")

	})
	t.Run("Exponential back off slows amount of Posts API calls if nothing is returned", func(t *testing.T) {
		// given
		client, recorder := GetRecorderClient(t, "posts")
		defer recorder.Stop()

		lc, _ := glaw.NewLemmyClient(url, "", jwt, client)

		// when
		lc.StreamNewPosts(2, nil)

		// let the endpoint listen for 30 seconds
		time.Sleep(30 * time.Second)

		// then
		assert.Greater(t, client.Transport.(*CountingTransport).Counter, 2, "Backoff was too slow, Counter should have hit Do() more than X times in given time")
		assert.Less(t, client.Transport.(*CountingTransport).Counter, 8, "Backoff didnt scale fast, Counter should have hit Do() less than X times in given time")
	})
}

func TestStreamComments(t *testing.T) {
	t.Run("Comments channel streams comments", func(t *testing.T) {
		//given
		client, recorder := GetRecorderClient(t, "comments")
		defer recorder.Stop()

		lc, _ := glaw.NewLemmyClient(url, "", jwt, client)

		stopChan := make(chan struct{})

		//when
		commentsChan := lc.StreamNewComments(5, stopChan)

		// pull posts out of the channel
		commentList := []glaw.Comment{}
		for i := 0; i < 10; i++ {
			p := <-commentsChan
			commentList = append(commentList, p)

		}

		// construct our expected slice
		// TODO: move this into a testing values file? or fixtures?
		var expectedComments []glaw.Comment
		inputString := `[{"id":181634,"creator_id":2086,"post_id":39509,"content":"test","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/181634","local":true,"path":"0.181634","distinguished":false,"language_id":0},{"id":171604,"creator_id":1537,"post_id":30391,"content":"@saltlake@voyager.lemmy.ml","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/171604","local":true,"path":"0.104129.104134.171604","distinguished":false,"language_id":37},{"id":171543,"creator_id":1537,"post_id":30391,"content":"cash or check?","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/171543","local":true,"path":"0.104129.104134.171543","distinguished":false,"language_id":37},{"id":163192,"creator_id":2086,"post_id":39204,"content":"Testing","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/163192","local":true,"path":"0.163192","distinguished":false,"language_id":0},{"id":163187,"creator_id":2086,"post_id":39509,"content":"test","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/163187","local":true,"path":"0.163187","distinguished":false,"language_id":0},{"id":158719,"creator_id":83995,"post_id":39509,"content":"Yy","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/158719","local":true,"path":"0.158719","distinguished":false,"language_id":0},{"id":104134,"creator_id":1528,"post_id":30391,"content":"$50 fee if you want to be un-banned","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/104134","local":true,"path":"0.104129.104134","distinguished":false,"language_id":37},{"id":104129,"creator_id":1528,"post_id":30391,"content":"banned","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/104129","local":true,"path":"0.104129","distinguished":true,"language_id":37},{"id":104052,"creator_id":1528,"post_id":3998,"content":"This comment is by a mod and is not distinguished","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/104052","local":true,"path":"0.104052","distinguished":false,"language_id":37},{"id":104050,"creator_id":1528,"post_id":3998,"content":"This comment is by a mod and is distinguished","removed":false,"deleted":false,"ap_id":"https://voyager.lemmy.ml/comment/104050","local":true,"path":"0.104050","distinguished":true,"language_id":37}]`
		err := json.Unmarshal([]byte(inputString), &expectedComments)
		if err != nil {
			t.Fatal(err)
		}

		//then
		assert.NotNil(t, commentList)
		assert.Equal(t, expectedComments, commentList)
	})
	t.Run("Handles closing channel correctly", func(t *testing.T) {
		// given
		client, recorder := GetRecorderClient(t, "comments")
		defer recorder.Stop()

		lc, err := glaw.NewLemmyClient(url, "", jwt, client)
		assert.Equal(t, nil, err)

		stopChan := make(chan struct{})

		//when
		commentsChan := lc.StreamNewComments(5, stopChan)

		time.Sleep(1 * time.Second)

		// close chan
		stopChan <- struct{}{}

		// pull all items out of channel
		for {
			_, ok := <-commentsChan
			if !ok {
				// Channel closed, and all values are received
				break
			}
		}

		// blocking check if channel is closed
		_, ok := <-commentsChan

		// then
		assert.False(t, ok, "Expected comments channel to be closed, but it's still open.")
	})
	t.Run("Exponential back off slows amount of Comments API calls if nothing is returned", func(t *testing.T) {
		// given
		client, recorder := GetRecorderClient(t, "comments")
		defer recorder.Stop()

		lc, _ := glaw.NewLemmyClient(url, "", jwt, client)

		// when
		lc.StreamNewComments(2, nil)

		// let the endpoint listen for 30 seconds
		time.Sleep(30 * time.Second)

		// then
		assert.Greater(t, client.Transport.(*CountingTransport).Counter, 2, "Backoff was too slow, Counter should have hit Do() more than X times in given time")
		assert.Less(t, client.Transport.(*CountingTransport).Counter, 8, "Backoff didnt scale fast, Counter should have hit Do() less than X times in given time")
	})
}

func TestSendPrivateMessage(t *testing.T) {
	t.Run("Can send a Private Message", func(t *testing.T) {
		t.Skip()
		//given
		client, recorder := GetRecorderClient(t, "send_private_message")
		defer recorder.Stop()
		lc, _ := glaw.NewLemmyClient("testURL", "token", "cookie", client)

		//when
		err := lc.SendPrivateMessage("Test send private message", 12345)

		//then
		assert.NotEqual(t, nil, lc)
		assert.Nil(t, err)
	})
}

func TestGetComment(t *testing.T) {
	t.Run("Can Get a Comment", func(t *testing.T) {
		t.Skip()
		//given
		client, recorder := GetRecorderClient(t, "get_comment")
		defer recorder.Stop()
		lc, err := glaw.NewLemmyClient("testURL", "token", "cookie", client)
		assert.Equal(t, nil, err)
		assert.NotNil(t, lc)

		//when
		comment, err := lc.GetComment(1234)

		//then
		assert.NotEqual(t, nil, comment)
		assert.Nil(t, err)
	})
	t.Run("Can handle comment error", func(t *testing.T) {
		t.Skip()
		//given
		client, recorder := GetRecorderClient(t, "get_comment")
		defer recorder.Stop()
		lc, err := glaw.NewLemmyClient("testURL", "token", "cookie", client)
		assert.Equal(t, nil, err)
		assert.NotNil(t, lc)

		//when
		comment, err := lc.GetComment(1234)

		//then
		assert.NotEqual(t, nil, comment)
		assert.Nil(t, err)
	})
}

func TestGetPost(t *testing.T) {
	t.Run("Can Get a Post", func(t *testing.T) {
		t.Skip()
		//given
		client, recorder := GetRecorderClient(t, "get_post")
		defer recorder.Stop()
		lc, err := glaw.NewLemmyClient("testURL", "token", "cookie", client)
		assert.Equal(t, nil, err)
		assert.NotNil(t, lc)

		//when
		post, err := lc.GetPost(1234)

		//then
		assert.NotEqual(t, nil, post)
		assert.Nil(t, err)
	})
	t.Run("Can handle post error", func(t *testing.T) {
		t.Skip()
		//given
		client, recorder := GetRecorderClient(t, "get_post")
		defer recorder.Stop()
		lc, err := glaw.NewLemmyClient("testURL", "token", "cookie", client)
		assert.Equal(t, nil, err)
		assert.NotNil(t, lc)

		//when
		post, err := lc.GetPost(1234)

		//then
		assert.NotEqual(t, nil, post)
		assert.Nil(t, err)
	})
}

func TestNewLemmyClient(t *testing.T) {
	t.Run("Can return a new client", func(t *testing.T) {
		//given
		client, recorder := GetRecorderClient(t, "posts")
		defer recorder.Stop()
		lc, err := glaw.NewLemmyClient("testURL", "token", "cookie", client)

		//when

		//then
		assert.NotEqual(t, nil, lc)
		assert.Nil(t, err)
	})

	t.Run("NewClient without base url returns nil client and error", func(t *testing.T) {
		//given
		client, recorder := GetRecorderClient(t, "posts")
		defer recorder.Stop()
		//when
		lc, err := glaw.NewLemmyClient("", "token", "cookie", client)
		//then
		assert.Nil(t, lc)
		assert.Equal(t, "url required", err.Error())
	})
}

func GetRecorderClient(t *testing.T, cassetteLocation string) (*http.Client, *recorder.Recorder) {
	r, err := recorder.New(fmt.Sprintf("fixture/%s_cassette", cassetteLocation))
	if err != nil {
		t.Fatal(err)
	}

	if r.Mode() != recorder.ModeRecordOnce {
		t.Fatal("Recorder should be in ModeRecordOnce")
	}
	// client := r.GetDefaultClient()

	client := r.GetDefaultClient()

	transport := &CountingTransport{
		Transport: r.GetDefaultClient().Transport,
	}

	client.Transport = transport
	return client, r
}

// custom struct so we can see how many times our client called the Do() function
type CountingTransport struct {
	Transport http.RoundTripper
	Counter   int
}

// Add a wrapper around the roundtrip to just increment the counter
func (c *CountingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.Counter++
	// fmt.Printf("Increased the counter. Counter: %v \n", c.Counter)
	return c.Transport.RoundTrip(req)
}
