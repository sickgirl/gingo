package ai

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/api"
	"github.com/songcser/gingo/pkg/response"
)

type Api struct {
	api.Api
	Service Service
}

func NewApi() Api {
	var app App
	s := NewService(app)
	baseApi := api.NewApi[App](s)
	return Api{Api: baseApi, Service: s}
}

func (a Api) Test(c *gin.Context) {
	str, err := a.Service.Test()

	if err != nil {
		fmt.Printf("请求错误: %s", err.Error())

		response.Fail(c)
		return
	}

	data := gin.H{
		"message": str,
	}
	response.OkWithData(data, c)
}
