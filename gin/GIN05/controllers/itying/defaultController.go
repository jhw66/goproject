package itying

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	//设置cookie
	//第三个参数为过期时间，<0表示删除cookie，>0表示添加cookie
	//第四个参数为cookie路径，第五个参数为cookie的Domain作用域
	//第六个参数为secure，为true表示在HTTP中无效，在HTTPS中才有效
	//第七个参数httpOnly是微软对cookie做的扩展，使通过程序将无法读取到cookie信息，防止XSS攻击产生
	c.SetCookie("username", "张三", 3600, "/", ".wjh.com", false, true)
	c.HTML(200, "itying.html", gin.H{
		"title": "首页-Index方法",
	})
}
func (con DefaultController) Get01(c *gin.Context) {
	//获取cookie
	username, _ := c.Cookie("username")
	c.String(200, "cookie="+username)
}
func (con DefaultController) Get02(c *gin.Context) {
	//获取cookie
	username, _ := c.Cookie("username")
	c.String(200, "cookie="+username)
}
func (con DefaultController) Delete(c *gin.Context) {
	c.SetCookie("username", "张三", -1, "/", ".wjh.com", false, true)
	c.String(200, "删除cookie")
}

func (con DefaultController) News(c *gin.Context) {
	c.String(200, "首页-News方法")
	username, _ := c.Get("username")
	if v, ok := username.(string); ok {
		c.String(200, v)
	} else {
		c.String(200, "获取用户失败")
	}
}

func (con DefaultController) About(c *gin.Context) {
	c.String(200, "首页-About方法")

	//在gorountine中要使用gin.Context要先复制一份使用
	cCp := c.Copy()
	c.String(200, c.GetString("username"))
	//注意此处不需要wg.Wait(),因为主程序一直在执行
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("DONE! in path" + cCp.Request.URL.Path)
	}()

}
