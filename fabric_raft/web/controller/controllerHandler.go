/**
功能测试

1. 接收表单



*/
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raft/fabric_raft/service"
)

type Application struct {
	Fabric *service.ServiceSetup
}

func Test(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "SUCCESS",
	})

}

func (app *Application) GetData(c *gin.Context) {

	key := c.DefaultQuery("key", "unknown")
	fmt.Println("查询..." + key)

	result, err := app.Fabric.FindInfo(key)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var m service.Mechanic
		json.Unmarshal(result, &m)
		fmt.Println("根据Key查询信息成功：")
		fmt.Println(m)
		c.JSON(http.StatusOK, m)
	}

}

func (app *Application) CreateData(c *gin.Context) {

	var m service.Mechanic

	if err := c.ShouldBind(&m); err != nil { // 不管是form,queryString,还是json，都自动判断接收
		fmt.Println(err.Error())
	} else {
		fmt.Println(m)
		fmt.Println("=======creating data========")

		msg, err := app.Fabric.Save(m)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("信息发布成功, 交易编号为: " + msg)
			c.JSON(http.StatusOK, gin.H{
				"transaction ID ": msg,
				"new key":         m.Key,
			})
		}

	}

}
