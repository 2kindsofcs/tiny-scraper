package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	endpoints := []string{"https://kubernetespodcast.com/"}

	// 5초에 한 번씩 웹페이지들을 긁어서 정보를 저장
	//
	for _, endpoint := range endpoints {
		content, err := scrape(endpoint)
		if err != nil {
			//do something
		}
		fmt.Println(content)
		fmt.Println("end")
		continue
	}
}

func scrape(endpoint string) (content string, err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	content = doc.Find("div.episode").First().Find("h3").Text()
	return content, nil
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
