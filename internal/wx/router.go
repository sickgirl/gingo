package wx

import (
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/middleware"
	"github.com/songcser/gingo/pkg/router"
)

func InitRouter(g *gin.RouterGroup) {

	gOpen := g.Group("wx_open")
	rO := router.NewRouter(gOpen)
	aO := NewApi()
	rO.BindGet("ask", aO.Ask) //询问

	gInner := g.Group("wx")
	gInner.Use(middleware.AuthMiddleware())

	r := router.NewRouter(gInner)
	a := NewApi()
	r.BindApi("", a)

	r.BindPost("test", a.Test)          //测试
	r.BindPost("bind_user", a.BindUser) //测试
	//r.BindGet("ask", a.Ask)                 //询问
	r.BindPost("send", a.Send)              //发送订阅消息
	r.BindGet("find_wx_user", a.FindWxUser) //发送订阅消息
	r.BindGet("ask_about", a.AskAbout)      //询问 + 向量化
}
