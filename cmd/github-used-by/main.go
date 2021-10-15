package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chyroc/go-ptr"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	githubToken := os.Getenv("GITHUB_TOKEN")

	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/trigger/:owner/:repo", func(c *gin.Context) {
		err := writeFile(githubToken, c.Param("owner"), c.Param("repo"))
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

func writeFile(githubToken, owner, repo string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	path := fmt.Sprintf("data/%s:_:_:%s", owner, repo)
	author := github.CommitAuthor{
		Name:  ptr.String("github-actions[bot]"),
		Email: ptr.String("41898282+github-actions[bot]@users.noreply.github.com"),
	}
	re, _, err := client.Repositories.UpdateFile(ctx, "chyroc", "github-used-by", path, &github.RepositoryContentFileOptions{
		Message:   ptr.String(fmt.Sprintf("add %s/%s", owner, repo)),
		Content:   []byte(owner + "/" + repo),
		Author:    &author,
		Committer: &author,
	})
	fmt.Println(re)
	return err
}
