package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chyroc/go-ptr"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/trigger-config/:owner/:repo", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		err := githubUpdateConfig(os.Getenv("GITHUB_TOKEN"), owner, repo)
		if err != nil {
			c.String(500, err.Error())
		} else {
			c.String(200, "success")
		}
	})
	r.Run("127.0.0.1:" + os.Getenv("SERVER_PORT"))
}

func githubUpdateConfig(token, owner, repo string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	reader, _, err := client.Repositories.DownloadContents(ctx, owner, repo, "config.yaml", nil)
	if err != nil {
		return err
	}

	defer reader.Close()
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	user := &github.CommitAuthor{
		Name:  ptr.String("github-actions[bot]"),
		Email: ptr.String("41898282+github-actions[bot]@users.noreply.github.com"),
	}
	newContent := fmt.Sprintf("%s- https://github.com/%s/%s\n", string(bs), owner, repo)
	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, "config.yaml", &github.RepositoryContentFileOptions{
		Message:   ptr.String(fmt.Sprintf("Update: %s/%s", owner, repo)),
		Content:   []byte(newContent),
		Author:    user,
		Committer: user,
	})
	return err
}
