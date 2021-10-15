package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chyroc/go-ptr"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	githubToken := os.Getenv("GITHUB_TOKEN")

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		MaxAge:       12 * time.Hour,
	}))
	r.Any("/trigger/:owner/:repo", func(c *gin.Context) {
		err := addRepo(githubToken, c.Param("owner"), c.Param("repo"))
		if err != nil {
			c.String(500, err.Error())
		} else {
			c.String(200, "ok")
		}
	})

	if err := r.Run("127.0.0.1:" + port); err != nil {
		panic(err)
	}
}

func addRepo(githubToken, owner, repo string) error {
	resp, _ := httpClient.Head(fmt.Sprintf("https://github.com/%s/%s", owner, repo))
	if resp != nil && resp.StatusCode == 404 {
		return nil
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	path := fmt.Sprintf("data/%s:_:_:%s", owner, repo)
	author := github.CommitAuthor{
		Name:  ptr.String("github-actions[bot]"),
		Email: ptr.String("41898282+github-actions[bot]@users.noreply.github.com"),
	}
	_, _, err := client.Repositories.UpdateFile(ctx, "chyroc", "github-used-by", path, &github.RepositoryContentFileOptions{
		Message:   ptr.String(fmt.Sprintf("add %s/%s", owner, repo)),
		Content:   []byte(owner + "/" + repo),
		Author:    &author,
		Committer: &author,
	})
	if err != nil {
		if strings.Contains(err.Error(), `"sha" wasn't supplied`) {
			return nil
		}
	}
	return err
}

var httpClient = &http.Client{Timeout: time.Second * 3}
