package main

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
)

func main() {
	srv := rest.MustNewServer(rest.RestConf{
		Port: 8888,
	})
	srv.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/hello",
		Handler: handle,
	})
	defer srv.Stop()

	// httpx.SetErrorHandler 仅在调用了 httpx.Error 处理响应时才有效。
	// 定义一个错误处理函数
	httpx.SetErrorHandler(func(err error) (int, any) {
		// 判断错误类型
		switch e := err.(type) {
		case *errors.CodeMsg:
			// 返回带有错误信息的 JSON 响应
			return http.StatusOK, xhttp.BaseResponse[struct{}]{
				Code: e.Code,
				Msg:  e.Msg,
			}
		default:
			// 返回服务器内部错误
			return http.StatusInternalServerError, nil
		}
	})

	srv.Start()
}

type HelloRequest struct {
	Name string `form:"name"`
}

type HelloResponse struct {
	Msg string `json:"msg"`
}

// 处理请求
func handle(w http.ResponseWriter, r *http.Request) {
	var req HelloRequest
	// 尝试解析请求数据
	if err := httpx.Parse(r, &req); err != nil {
		httpx.Error(w, err) // 触发错误处理
		fmt.Println("Error parsing request")
		return
	}

	// 模拟业务逻辑错误
	if req.Name == "error" {
		httpx.Error(w, errors.New(400, "参数错误"))
		fmt.Println("Simulated error triggered")
		return
	}

	// 返回正常响应
	httpx.OkJson(w, HelloResponse{
		Msg: "hello " + req.Name,
	})
	fmt.Println("Response sent")
}
