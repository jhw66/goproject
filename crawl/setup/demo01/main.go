package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	url := "https://book.douban.com/"
	// 功能：定义一个变量 url，保存要请求的网页地址。
	// 类型：string
	// 作用：用于后续 HTTP 请求。

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 	函数解释：
	// http.NewRequest(method string, url string, body io.Reader) (*http.Request, error)
	//      method：请求方法，例如 "GET"、"POST"。
	//      url：请求目标 URL。
	//      body：请求体，对于 GET 请求通常是 nil。
	// 返回值：
	//      req：*http.Request 类型，表示 HTTP 请求对象。
	//      err：error 类型，如果请求创建失败，会返回错误信息。
	// 作用：
	//      生成一个请求对象，后续可以设置请求头、Cookie 等。
	// 特性：
	//      可以复用 *http.Request 对象。
	//      并未发送请求，仅构造请求。

	//-------------------------------------------------------------------------------------------

	//操作 HTTP 请求头（Headers）
	// 	req.Header.Set 函数
	// 函数签名
	// func (h Header) Set(key, value string)
	// 所属类型：Header，属于 http 包中的 http.Request 的 Header 字段。
	// 作用：
	// 在 HTTP 请求头中设置（或覆盖）某个字段的值，如果同名字段已经存在，会覆盖原值。
	// 参数：
	// key：HTTP 头字段名（string），区分大小写（通常首字母大写，如 "User-Agent"）。
	// value：HTTP 头字段值（string）。
	// 返回值：
	// 无返回值（void），只是修改了请求头的内部映射。
	//
	// Header 类型说明：
	// type Header map[string][]string
	// 内部是 map[string][]string，支持一个字段有多个值：
	// Set → 设置单个值，会覆盖旧值。
	// Add → 添加一个值，不覆盖已有值。
	// 常用方法：
	// Set(key, value string) → 设置/覆盖
	// Add(key, value string) → 添加
	// Get(key string) string → 获取第一个值
	// Del(key string) → 删除字段
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	// 作用：告诉服务器客户端是什么浏览器/操作系统。
	// 参数说明：
	// 字符串一般模拟浏览器的标识。
	// 这里模拟的是 Chrome 120、Windows 10。
	// 为什么要设置：
	// 很多网站会根据 User-Agent 返回不同的页面，或者防止非浏览器访问。
	// 可替换参数：
	// 任意浏览器或设备标识字符串，如：
	// "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)..." → 手机浏览器
	// "curl/7.68.0" → curl 命令
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	// 	作用：告诉服务器客户端 能接受的内容类型（MIME Type）。
	// 参数说明：
	// "text/html" → HTML 页面
	// "application/xhtml+xml" → XHTML 页面
	// "application/xml;q=0.9" → XML 页面，q=0.9 表示优先级 0.9
	// "*/*;q=0.8" → 其它任意类型，优先级 0.8
	// q 参数：
	// 用来表示 优先级/权重，取值 0~1，越大越优先
	// 可设置参数：
	// 根据请求需求，可设置：
	// "application/json"
	// "image/webp,image/apng,*/*;q=0.8"
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//作用：告诉服务器 客户端希望返回的语言。
	// 参数说明：
	// "zh-CN" → 简体中文（中国）
	// "zh;q=0.9" → 中文其他变体，优先级 0.9
	// 可设置参数：
	// "en-US,en;q=0.9" → 美式英语
	// "fr-FR,fr;q=0.8" → 法语
	req.Header.Set("Connection", "keep-alive")
	// 	作用：告诉服务器是否复用 TCP 连接。
	// 参数说明：
	// "keep-alive" → 保持连接，多个请求可复用同一个 TCP 连接
	// "close" → 请求完成后关闭连接
	// 可设置参数：
	// "keep-alive"（复用连接，提高性能）
	// "close"（每次请求新建连接，安全但慢）

	//-------------------------------------------------------------------------------------------

	client := &http.Client{}
	// 	功能：创建 HTTP 客户端对象。
	// 类型：*http.Client。
	// 作用：
	// 用来发送 HTTP 请求。
	// 可以设置超时、重定向策略、代理等。
	// 特性：
	// 默认不设置超时，可能导致请求无限挂起（生产环境建议设置 Timeout）。
	// 可以复用 *http.Client 对象，提高连接效率。

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// 	功能：发送 HTTP 请求并获取响应。
	// 函数解释：
	// client.Do(req *http.Request) (*http.Response, error)
	// 返回值：
	// resp：*http.Response，包含状态码、头、Body 等响应信息。
	// err：error，请求失败会返回错误。
	// 作用：
	// 真正向服务器发起请求。
	// 网络层、DNS 解析、TLS 握手都会在这里发生。
	// 特性：
	// 阻塞调用，直到收到响应或超时。
	// resp.Body 是 io.ReadCloser，需要关闭释放资源。

	defer resp.Body.Close()
	// 	功能：延迟关闭响应体。
	// 类型：resp.Body 是 io.ReadCloser。
	// 作用：
	// 确保网络连接被释放，防止资源泄漏。
	// 特性：
	// defer 在函数返回时执行。
	// 对长连接尤其重要，否则可能耗尽系统连接。

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code: %d\n", resp.StatusCode)
		return
	}
	// 	功能：检查响应状态码是否为 200。
	// 类型：
	// resp.StatusCode：int，HTTP 状态码。
	// http.StatusOK：常量 200。
	// 作用：
	// 只有返回 200 才认为请求成功。
	// 否则打印错误状态码并提前退出。
	// 特性：
	// 可以根据不同状态码处理不同逻辑（如 3xx 重定向、4xx/5xx 错误）。

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// 	功能：读取响应体的全部内容。
	// 函数解释：
	// io.ReadAll(r io.Reader) ([]byte, error)
	// 参数 r 是实现了 io.Reader 接口的对象，这里是 resp.Body。
	// 返回值：
	// body：[]byte，包含网页内容（HTML）。
	// err：读取失败时返回错误。
	// 作用：
	// 将响应体完全加载到内存。
	// 特性：
	// 对大文件可能消耗大量内存。
	// 适合小网页抓取。
	// []byte 可以通过 string(body) 转成字符串。

	//-------------------------------------------------------------------------------------------

	file, err := os.OpenFile("../test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 写入文件
	Writer := bufio.NewWriter(file)
	Writer.WriteString(string(body))
	Writer.Flush()

	fmt.Println("输出已写入 output.txt")
	fmt.Println(len(body))
	fmt.Println(resp.Status)
}

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// func main() {
// 	url := "https://book.douban.com/"

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Status:", resp.Status)
// 	fmt.Println("Body:", string(body))
// }

// 返回状态码 418
// 418 I'm a teapot 是一个“彩蛋”状态码，原本是 RFC 2324（愚人节协议）里定义的，意思是“我是一台茶壶，不能泡咖啡”。
// 很多网站在检测到异常请求（比如没有浏览器标识的爬虫请求）时，会返回这个状态码作为防爬手段。
// 也就是说，豆瓣检测到你的 http.Get 请求，觉得它是“机器人”，于是返回了 418。
