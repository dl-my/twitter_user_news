package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"twitter_user_news/api"
)

func Run(port int) {
	r := gin.Default()

	// 注册路由
	registerRoutes(r)

	// 启动服务
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *gin.Engine) {
	r.POST("/add", api.Add)
	r.POST("/del", api.Del)
	r.POST("/del_all", api.DelAll)
	r.GET("/list", api.List)
	r.POST("/reload_cookie", api.ReloadCookie)
}
