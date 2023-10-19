package ai

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/router"
)

func InitRouter(g *gin.RouterGroup) {
	r := router.NewRouter(g.Group("ai"))
	a := NewApi()
	r.BindApi("", a)
	r.BindGet("test", a.Test)          //测试
	r.BindGet("ask", a.Ask)            //询问
	r.BindGet("ask_about", a.AskAbout) //询问 + 向量化
}
