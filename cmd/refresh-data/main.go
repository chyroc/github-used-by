package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	repos, err := load()
	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		count, err := getCount(repo)
		if err != nil {
			panic(err)
		}
		if count == "" {
			continue
		}
		badge, err := getBadge(fmt.Sprintf("Used by-%s-brightgreen", count))
		if err != nil {
			panic(err)
		}
		file := "docs/" + repo + ".svg"
		if err = os.MkdirAll(filepath.Dir(file), 0o777); err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile(file, []byte(badge), 0o666); err != nil {
			panic(err)
		}
	}
}

func load() (res []string, err error) {
	fs, err := ioutil.ReadDir("data")
	if err != nil {
		return nil, err
	}
	for _, v := range fs {
		res = append(res, strings.ReplaceAll(v.Name(), ":_:_:", "/"))
	}
	return res, nil
}

func getCount(repo string) (string, error) {
	res, err := httpClient.Get(fmt.Sprintf("https://github.com/%s", repo))
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	match := regCount.FindStringSubmatch(string(bs))
	if len(match) != 2 {
		return "", nil
	}
	return match[1], nil
}

func getBadge(key string) (string, error) {
	res, err := httpClient.Get("https://img.shields.io/badge/" + url.PathEscape(key))
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

var (
	regCount   = regexp.MustCompile(`Used by.*?<span.*?>(\d)+</span>`)
	httpClient = &http.Client{Timeout: time.Second * 3}
)
