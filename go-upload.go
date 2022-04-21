package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/zh-five/xdaemon"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

/// 这里是用embed把静态文件打包入二进制里面
//go:embed templates/*
var f embed.FS

var (
	port   string
	daemon bool
)

/// 初始化参数默认值
func init() {
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", "80", "监听端口号")
	serverCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "是否为守护进程模式")
}

/// 定义 -h 的提示信息
var serverCmd = &cobra.Command{
	Use:     "",
	Short:   "gin实现简版文件服务",
	Example: "go-upload -p 8080 -d",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

/// 上传方法
func upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	filename := header.Filename
	out, err := os.Create("public/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "up.html", gin.H{})
}

/// 主方法
func run() {
	log.Println("使用 go-upload -h 查看更多命令")
	/// 判断创建上传目录
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/public/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
		}
	}

	/// 后台运行
	if daemon == true {
		d := xdaemon.NewDaemon(logFilePath + "go-upload.log")
		d.MaxCount = 10
		d.Run()
	}

	router := gin.Default()
	/// 装入静态文件
	templ := template.Must(template.New("").ParseFS(f, "templates/*"))
	router.SetHTMLTemplate(templ)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.StaticFileFS("/favicon.ico", "templates/favicon.ico", http.FS(f))

	router.POST("/upload", upload)

	router.StaticFS("/files", gin.Dir("public", true))

	router.Run("0.0.0.0:" + port)

}

/// 入口
func main() {
	if err := serverCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
