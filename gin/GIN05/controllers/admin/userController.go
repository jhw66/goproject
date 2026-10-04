package admin

import (
	"fmt"
	"mygin/GIN05/models"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {
	c.String(200, "用户列表")
	con.success(c)
	con.fail(c)
	username, _ := c.Get("username")
	if v, ok := username.(string); ok {
		c.String(200, v)
	} else {
		c.String(200, "获取用户失败")
	}
}

func (con UserController) Addfiles(c *gin.Context) {
	c.HTML(200, "addfiles.html", gin.H{})
}

func (con UserController) DoUpload(c *gin.Context) {
	// username := c.PostForm("username")
	// file1, err := c.FormFile("face1")
	// dst1 := path.Join("./static/upfiles", file1.Filename)
	// //.代表的是GIN05目录下，因为其实这都是main.go中函数的调用
	// if err == nil {
	// 	c.SaveUploadedFile(file1, dst1)
	// 	c.JSON(200, gin.H{
	// 		"username": username,
	// 		"success":  true,
	// 		"dst":      dst1,
	// 	})
	// } else {
	// 	c.String(200, "文件上传失败")
	// }
	// file2, err := c.FormFile("face2")
	// dst2 := path.Join("./static/upfiles", file2.Filename)

	// if err == nil {
	// 	c.SaveUploadedFile(file2, dst2)
	// 	c.JSON(200, gin.H{
	// 		"username": username,
	// 		"success":  true,
	// 		"dst":      dst2,
	// 	})
	// } else {
	// 	c.String(200, "文件上传失败")
	// }

	username := c.PostForm("username")

	//1，获取上传的文件
	form, _ := c.MultipartForm()
	files := form.File["face[]"]

	for _, file := range files {
		//2，获取后缀名，判断类型是否正确
		extName := path.Ext(file.Filename)
		allowExtMap := map[string]bool{
			".jpg":  true,
			".png":  true,
			".gif":  true,
			".jpeg": true,
		}
		if _, ok := allowExtMap[extName]; !ok {
			c.String(200, "上传文件类型不合法\n")
			continue
		}
		//3，创建图片保存目录
		day := models.GetDate()
		dir := path.Join("./static/upload/", day)

		// os.Mkdir = "创建这个目录"（要求父目录已存在）
		// os.MkdirAll = "创建这个目录和它所有不存在的父目录"
		err := os.MkdirAll(dir, 0666) //0666表示所有人都有写和读的权限
		if err != nil {
			fmt.Println(err)
			c.String(200, "MkdirAll失败")
			return
		}

		//4，生成文件名称和文件保存目录
		fileName := strconv.FormatInt(models.GetUnix(), 10) + extName
		dst := dir + "/" + fileName

		//5，执行上传
		c.SaveUploadedFile(file, dst)
		c.String(200, username+":"+file.Filename+" 保存成功\n")

		time.Sleep(time.Second)
	}
}

func (con UserController) Edit(c *gin.Context) {
	c.String(200, "用户编辑--edit")
}
