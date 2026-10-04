package main

import (
	"fmt"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Title string `json:"title"`
	News  string `json:"news"`
}

func UnixToTime(timestamp int64) string {
	// fmt.Println(timestamp)
	//这个打印是打印在控制台的，不是网页上
	//需要把时间戳转换成时间格式
	//2006-01-02 15:04:05
	//time.Unix()第一个参数是秒，第二个参数是纳秒
	t := time.Unix(timestamp, 0)
	return t.Format("2006-1-02 15:04:05")
}
func Println(arge ...interface{}) interface{} {
	var result interface{}
	for k, v := range arge {
		if k == 0 {
			result = v
			continue
		}
		result = fmt.Sprintf("%v %v", result, v)
	}
	return result
}

func main() {
	r := gin.Default()
	//自定义模板函数
	//放在加载模板文件的前面
	r.SetFuncMap(template.FuncMap{
		"Len": func(s string) int {
			return len(s)
		},
		"UnixToTime": UnixToTime,
		"Println":    Println,
	})
	//放在路由器初始化的后面，路由设置的前面
	//加载模板文件，**表示可以加载多级目录
	//模板文件放在项目的template目录下
	r.LoadHTMLGlob("template/**/*")
	a := &Article{
		Title: "新闻标题",
		News:  "这是新闻内容",
	}

	// 	r.Static 的作用是：
	// 把一个 URL 路径前缀，映射到服务器本地的一个目录，用来对外暴露静态资源
	// 也可以理解为：
	// “当浏览器访问某个路径时，Gin 去某个文件夹里找文件并返回”
	//配置静态web目录，第一个参数表示路由，第二个参数表示映射的目录
	r.Static("/static", "./static")
	//这样子写就在http://localhost:8080/static/CSS/base.css
	//r.Static("/sss", "./static")
	//这样子写就在http://localhost:8080/sss/CSS/base.css

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title":      a.Title,
			"timestamp":  time.Now().Unix(),
			"content":    []string{"hello", "gin", "template"},
			"mixcontent": []interface{}{"hello", "wjh", "6666", "again"},
		})
	})
	r.GET("/news", func(c *gin.Context) {
		c.HTML(200, "news.html", gin.H{
			"title": a.Title,
			"news":  a.News,
			"newsList": []interface{}{
				&Article{
					Title: "新闻标题1",
					News:  "这是新闻内容1",
				},
				&Article{
					Title: "新闻标题2",
					News:  "这是新闻内容",
				},
			},
		})
	})
	r.GET("/admin/index", func(c *gin.Context) {
		c.HTML(200, "admin/index.html", gin.H{
			"title":  "后台首页",
			"scores": 88,
		})
	})
	r.GET("/admin/news", func(c *gin.Context) {
		c.HTML(200, "admin/news.html", gin.H{
			"news": "后台新闻内容",
			"newsList": []interface{}{
				&Article{
					Title: "后台新闻标题1",
					News:  "这是后台新闻内容1",
				},
				&Article{
					Title: "后台新闻标题2",
					News:  "这是后台新闻内容2",
				},
				gin.H{
					"Title": "这是后台新闻内容3",
					"News":  "这是后台新闻内容3",
				},
				map[string]string{
					"Title": "这是后台新闻内容4",
					"News":  "这是后台新闻内容4",
				},
			},
		})
	})

	r.Run(":8080")
}
