# glaw
Golang Lemmy API Wrapper

# download
```go get github.com/adadkins/glaw```

# example usage
```
    // make a lemmy client, pass in either your Cookie JWT or Auth JWT depending on how your instance auths
    client := http.Client{}
	lc, _ := glaw.NewLemmyClient(url, "Auth JWT", "cookieJWT", client, nil)

    // get a comment
    comment, _ := lc.GetComment(1)
    fmt.Println(comment)
```