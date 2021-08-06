/**
  author: kevin
*/
package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raft/fabric_raft/web/controller"
)

func WebStart(app *controller.Application) {

	r := gin.Default()
	r.Static("/static", "static")

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "SUCCESS",
		})
	})

	r.GET("/fab", app.GetData)
	r.POST("/fab", app.CreateData)

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
