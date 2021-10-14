package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	//go:embed config.yaml
	config   string
	httpCli  = &http.Client{Timeout: time.Second * 3}
	regCount = regexp.MustCompile(`Used by.*?<span .*?>(\d)+</span>`)
)


func main() {
	err := f()
	if err != nil {
		log.Fatalln(err)
	}
}

func f() error {
	repos := []string{}
	err := yaml.Unmarshal([]byte(config), &repos)
	if err != nil {
		return err
	}
	var errs []string
	for _, repo := range repos {
		html, err := getHtml(repo)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		count, err := getCount(html)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Println(getRepoName(repo), count)
	}
	fmt.Println(repos)

	return nil
}

func getRepoName(uri string) string {
	uri = strings.TrimRight(uri, "/")
	uris := strings.Split(uri, "/")
	if len(uris) < 2 {
		return uri
	}

	return strings.Join(uris[len(uris)-2:], "/")
}

func getHtml(uri string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	resp, err := httpCli.Do(req)
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func getCount(html string) (int, error) {
	res := regCount.FindStringSubmatch(html)
	if len(res) != 2 {
		return 0, nil
	}
	i, _ := strconv.ParseInt(res[1], 10, 64)
	return int(i), nil
}
