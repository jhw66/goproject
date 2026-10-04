package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func NewRequest(method string, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	return req, err
}
func main() {
	req, err := NewRequest("GET", "https://news.sohu.com")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	count := 0
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		href := s.AttrOr("href", "")

		if !strings.Contains(href, "/a/") || title == "" {
			return
		}

		count++ // 自己计数
		fmt.Printf("[%d %d] %s\n链接: %s\n\n", count, i+1, title, href)
	})

}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/PuerkitoBio/goquery"
// )

// func NewRequest(method, url string) (*http.Request, error) {
// 	req, err := http.NewRequest(method, url, nil)
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
// 	return req, err
// }

// func main() {
// 	url := "https://www.sohu.com/a/956211860_120094090"

// 	req, err := NewRequest("GET", url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	client := &http.Client{Timeout: 10 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// === 1. 标题 ===
// 	title := strings.TrimSpace(doc.Find("#article-container > div.left.main > div:nth-child(1) > div > div.text-title > h1").Text())

// 	// === 2. 发布时间 ===
// 	pubTime := strings.TrimSpace(doc.Find("#news-time").Text())

// 	// === 3. 正文内容 ===
// 	content := ""
// 	doc.Find("#mp-editor p").Each(func(i int, s *goquery.Selection) {
// 		text := s.Text()
// 		if len(text) > 0 {
// 			content += text + "\n"
// 		}
// 	})

// 	// 打印结果
// 	fmt.Println("标题：", title)
// 	fmt.Println("时间：", pubTime)
// 	fmt.Println("正文：\n", content)
// }
