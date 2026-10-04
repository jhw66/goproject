// demo07：net/http 与 context 的集成（最常用场景）。
//
//	服务端：从 r.Context() 感知客户端断开 / 超时，尽早停止工作。
//	客户端：用 http.NewRequestWithContext 控制单次请求的超时/取消。
//
// 本例用 httptest.NewServer 启动一个本地 Server，不会挂起。
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("[server] 收到请求，开始处理 (2s)")

	select {
	case <-time.After(2 * time.Second):
		fmt.Fprintln(w, "ok")
		fmt.Println("[server] 处理完毕")
	case <-ctx.Done():
		fmt.Println("[server] 客户端提前结束:", ctx.Err())
		fmt.Fprintln(w, ctx.Err())
	}
}

func main() {
	ts := httptest.NewServer(http.HandlerFunc(slowHandler))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[client] 请求失败:", err)
		fmt.Println("[client] ctx.Err():", ctx.Err())
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("[client] 收到:", string(body))
	}

	time.Sleep(200 * time.Millisecond)
}
