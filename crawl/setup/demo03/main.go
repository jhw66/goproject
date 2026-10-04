package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
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
	req, err := NewRequest("GET", "https://news.sohu.com/")
	if err != nil {
		log.Fatal(err)
	}

	//请求超时（timeout） 指的是：
	// 当程序发出一个 HTTP 请求后，如果在一定时间内没有收到服务器响应，就自动终止请求。

	// client := &http.Client{
	// 	Timeout: 5 * time.Second, // 最多等 5 秒
	// }
	// // 	这个 Timeout 包括：
	// // DNS 解析
	// // TCP 连接
	// // 发送请求
	// // 读取响应
	// // 整个过程总时间不超过 5 秒。

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second, // 连接超时
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ResponseHeaderTimeout: 5 * time.Second, // 读取响应头超时
		},
	}
	// 	Client.Timeout	整个请求全过程	简单爬虫、通用请求
	// Dialer.Timeout	TCP 连接时间	网络慢或不稳定时
	// ResponseHeaderTimeout	等待响应头时间	服务端响应慢时
	// context.WithTimeout()	单独控制每个请求	灵活控制单个请求

	resp, _ := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	file, _ := os.OpenFile("./test.txt", os.O_WRONLY|os.O_CREATE, 0644)
	Writer := bufio.NewWriter(file)
	_, err = Writer.WriteString(string(body))
	if err != nil {
		log.Println("WriteString 错误:", err)
	}
	err = Writer.Flush()
	if err != nil {
		log.Println("Flush 错误:", err)
	}
	fmt.Println("写入成功")
	defer resp.Body.Close()
}
