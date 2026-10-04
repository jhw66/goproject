package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiController struct{}

func (api ApiController) Index(c *gin.Context) {
	c.String(200, "API控制器-Index方法"+" "+c.GetString("username"))
	//没有使用中间件所以得不到中间件的username:"张三"
}

func (api ApiController) Add(c *gin.Context) {
	time.Sleep(time.Second)
	c.String(200, "API控制器-Add方法"+" "+c.GetString("username"))
	c.Set("username", "回")
}
func (api ApiController) Add2(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.String(200, c.GetString("username"))
}

func (api ApiController) Edit(c *gin.Context) {
	time.Sleep(time.Second)
	c.String(200, "API控制器-Edit方法\n")
}

func Printaaa(c *gin.Context) {
	fmt.Println("aaa")
	//abort表示终止其他处理程序
	c.Abort()
	fmt.Println("aaa")

}

func Timelong(c *gin.Context) {
	begin := time.Now().UnixNano()
	// next表示调用该请求的剩余处理程序
	c.Next()
	end := time.Now().UnixNano()
	long := end - begin
	fmt.Println(long)
	c.String(200, strconv.Itoa(int(long)))
}
