package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	endpoints := []string{"https://kubernetespodcast.com/"}

	// 5초에 한 번씩 웹페이지들을 긁어서 정보를 저장
	//
	for _, endpoint := range endpoints {
		title, link, err := scrapeKubePodcast(endpoint)
		if err != nil {
			//do something
		}
		fmt.Println(title, link)
		scrapeHerokuPodcast("https://www.heroku.com/podcasts/codeish")
		fmt.Println("end")
		continue
	}
}

func scrapeKubePodcast(endpoint string) (title, link string, err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	doc.Find("div.episode h3").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(strings.ToLower(s.Text()), "istio") {
			l, ok := s.Find("a").First().Attr("href")
			if !ok {
				panic("link should exist.")
			}
			title = s.Text()
			link = l
		}
	})

	return title, link, nil
}

func scrapeHerokuPodcast(endpoint string) (title, link string, err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	doc.Find("div.episode-text-summary h2").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(strings.ToLower(s.Text()), "engineering") {
			l, ok := s.Find("a").First().Attr("href")
			if !ok {
				panic("link should exist.")
			}
			title = s.Text()
			link = "https://heroku.com" + l
		}
	})

	return title, link, nil
}

// utf-8인지 검사하는 유틸 함수. 꼭 여기 있을 필요는 없음.
func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}
