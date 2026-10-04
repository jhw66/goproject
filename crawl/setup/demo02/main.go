// 完全可以不用每次都手动写这么多 Header，Go 里有几种方式和第三方库可以帮你自动处理或简化这些操作。我们来详细说一下：

// 1️⃣ 使用 http.Client + http.Transport 自动处理部分 Header
// Go 自带的 http.Client 对一些常用 Header（如 Connection: keep-alive、User-Agent）并不会自动添加，
// 但你可以写一个 自定义 Transport 或 封装函数，在每次请求前统一设置默认 Header。

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// )

// func NewRequest(method string, url string) (*http.Request, error) {
// 	req, err := http.NewRequest(method, url, nil)
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
// 		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
// 	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
// 	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
// 	req.Header.Set("Connection", "keep-alive")
// 	return req, err
// }
// func main() {
// 	url := "https://book.douban.com/"
// 	req, err := NewRequest("GET", url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	Client := new(http.Client)
// 	resp, _ := Client.Do(req)
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		log.Fatal(resp.StatusCode)
// 	}
// 	body, _ := io.ReadAll(resp.Body)
// 	file, _ := os.OpenFile("./test.txt", os.O_CREATE|os.O_WRONLY, 0644)
// 	Writer := bufio.NewWriter(file)
// 	Writer.WriteString(string(body))
// 	Writer.Flush()
// 	fmt.Println(len(body))
// }

// 2️⃣ 使用第三方库 colly / gocolly（爬虫库）
// Colly是 Go 的一个非常流行的爬虫库，它可以帮你 自动处理常用浏览器 Header、Cookie、代理。

package main

import (
	"bufio"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	file, _ := os.OpenFile("./test02.txt", os.O_CREATE|os.O_WRONLY, 0644)
	Writer := bufio.NewWriter(file)
	// c.OnHTML("a[href]", func(h *colly.HTMLElement) {
	// 	h.Request.Visit(h.Attr("href"))
	// })
	c.OnResponse(func(r *colly.Response) {
		Writer.Write(r.Body)
	})
	Writer.Flush()
	c.Visit("https://book.douban.com/")
}

// 这里我们使用了 Colly 库的 OnResponse 方法来处理响应。
// 当 Colly 收到响应时，它会调用我们在 OnResponse 中定义的回调函数。
// 在这个回调函数里，我们可以直接处理响应数据，比如写入文件。
// 虽然 c.Visit() 写在 c.OnResponse() 后面，但执行顺序仍然是 Visit 先发请求、OnResponse 后执行。
// 也就是说：
// c.OnResponse(...) // 注册回调函数（只是定义好“等下要做什么”）
// c.Visit(...)      // 发送请求 -> 收到响应 -> 调用回调函数
// 核心原理：事件回调机制（Event-driven）
// Colly（和很多现代框架一样）是事件驱动模型（Event-driven）。
// 这意味着：
// 你在 OnResponse 里写的函数并不是马上执行的；
// 它是一个回调函数（Callback）；
// Colly 会在适当的时机触发它。

//3️⃣ 使用 http.Request 的 Header.Add 与 map 封装
// 如果不想用第三方库，也可以自己封装一个 map，统一管理默认 Header：

// defaultHeaders := map[string]string{
//     "User-Agent":      "MyGoClient/1.0",
//     "Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
//     "Accept-Language": "zh-CN,zh;q=0.9",
//     "Connection":      "keep-alive",
// }

// req, _ := http.NewRequest("GET", url, nil)
// for k, v := range defaultHeaders {
//     req.Header.Set(k, v)
// }
