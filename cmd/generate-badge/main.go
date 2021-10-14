package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var httpCli = &http.Client{Timeout: time.Second * 3}

func main() {
	err := f()
	if err != nil {
		panic(err)
	}
}

func f() error {
	counts, err := readRepoCount()
	if err != nil {
		return err
	}
	errs := []string{}
	for repo, count := range counts {
		badge, err := generateBadgeSvg(count)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		file := repo + ".svg"
		if err = os.MkdirAll(filepath.Dir(file), 0777); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if err = ioutil.WriteFile(file, []byte(badge), 0666); err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, " / "))
	}

	return nil
}

func readRepoCount() (map[string]int, error) {
	res := map[string]int{}
	filepath.Walk("data", func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".txt") {
			repo := path[len("data/") : len(path)-len(".txt")]
			bs, _ := ioutil.ReadFile(path)
			count, _ := strconv.ParseInt(string(bs), 10, 64)
			res[repo] = int(count)
		}
		return nil
	})
	return res, nil
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

func generateBadgeSvg(count int) (string, error) {
	uri := fmt.Sprintf("https://img.shields.io/badge/Used%%20by-%d-brightgreen", count)
	return getHtml(uri)
}
