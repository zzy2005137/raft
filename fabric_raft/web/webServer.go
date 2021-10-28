/**
  author: kevin
*/
package web

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/raft/fabric_raft/web/controller"
)

func WebStart(app *controller.Application) {

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	})

	r.Static("/static", "static")
	//使用绝对路径
	r.LoadHTMLGlob(filepath.Join(os.Getenv("GOPATH"), "src/github.com/raft/fabric_raft/web/tpl/*"))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil) //HTML渲染

	})

	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "SUCCESS",
		})
	})

	r.GET("/detail", app.GetDetails)

	r.GET("/fab", app.GetData)
	r.POST("/fab", app.CreateData)
	// r.POST("/post", app.Test)

	r.Run(":8080")

	// fs := http.FileServer(http.Dir("web/static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	// http.HandleFunc("/", app.IndexView)
	// http.HandleFunc("/index.html", app.IndexView)
	// http.HandleFunc("/setInfo.html", app.SetInfoView)
	// http.HandleFunc("/setReq", app.SetInfo)
	// http.HandleFunc("/queryReq", app.QueryInfo)

	// fmt.Println("启动Web服务, 监听端口号: 9000")
	// err := http.ListenAndServe(":9000", nil)
	// if err != nil {
	// 	fmt.Println("启动Web服务错误")
	// }

}
