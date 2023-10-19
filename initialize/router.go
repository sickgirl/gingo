package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/internal/ai"
	"github.com/songcser/gingo/internal/app"
	"github.com/songcser/gingo/middleware"
	"github.com/songcser/gingo/utils"
	"net/http"
)

func HealthCheck(g *gin.Context) {
	data := gin.H{
		"message": "everyThing is Ok!",
	}
	g.JSON(http.StatusOK, data)
}

func Routers() *gin.Engine {

	if err := utils.Translator("zh"); err != nil {
		config.GVA_LOG.Error(err.Error())
		return nil
	}

	Router := gin.Default()
	//gin.SetMode(gin.DebugMode)

	Router.Use(middleware.Recovery())
	Router.Use(middleware.Logger())
	HealthGroup := Router.Group("/health")
	{
		// 健康监测
		HealthGroup.GET("/hhh", HealthCheck)
	}

	ApiGroup := Router.Group("api/v1")
	app.InitRouter(ApiGroup)

	//注册ai 相关路由
	ai.InitRouter(ApiGroup)

	return Router
}
