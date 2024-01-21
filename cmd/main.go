package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var wg sync.WaitGroup
	resChan := make(chan Result, 2)
	results := make([]Result, 0, 2)

	c := time.Tick(3 * time.Second)

	for {
		select {
		case <-c:
			wg.Add(2)
			go scrapeKubePodcast("https://kubernetespodcast.com/", &wg, resChan)
			go scrapeHerokuPodcast("https://www.heroku.com/podcasts/codeish", &wg, resChan)
			wg.Wait()

		case data := <-resChan:
			results = append(results, data)
			if len(results) == 2 {
				fmt.Println(results)
				//TODO: DB에 저장
				results = make([]Result, 0, 2)
			}
		}
	}
}

type Result struct {
	title string
	link  string
}

func scrapeKubePodcast(endpoint string, wg *sync.WaitGroup, resChan chan<- Result) {
	defer wg.Done()
	var title, link string

	resp, err := http.Get(endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
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
	resChan <- Result{
		title: title,
		link:  link,
	}
	return
}

func scrapeHerokuPodcast(endpoint string, wg *sync.WaitGroup, resChan chan<- Result) {
	defer wg.Done()
	var title, link string

	resp, err := http.Get(endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
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

	resChan <- Result{
		title: title,
		link:  link,
	}
	return
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
