package wx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/songcser/gingo/pkg/api"
	"github.com/songcser/gingo/pkg/response"
	"log"
	"strconv"
)

type Api struct {
	api.Api
	Service Service
}

func NewApi() Api {
	var app TWxUsers
	s := NewService(app)
	baseApi := api.NewApi[TWxUsers](s)
	return Api{Api: baseApi, Service: s}
}

func (a Api) Test(c *gin.Context) {

	body := BindRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("参数错误: %s", err.Error())
		response.Fail(c)
		return
	}

	str, err := a.Service.Test("")

	if err != nil {
		fmt.Printf("请求错误: %s", err.Error())

		response.Fail(c)
		return
	}
	userID, _ := c.Get("user_id")

	// 将 userID 转换为 string 类型
	intUserID, ok := userID.(int)
	if !ok {
		// 转换失败，处理错误
		log.Println("Failed to convert user_id to string")
		// 进行错误处理逻辑...
		return
	}

	strUserID := strconv.Itoa(intUserID)

	// 在这里使用转换后的 string 类型的 userID 进行后续操作
	fmt.Println("User ID:", strUserID)

	data := gin.H{
		"message": str + body.Phone + "当前用户id:" + strUserID,
	}
	response.OkWithData(data, c)
}

func (a Api) BindUser(c *gin.Context) {

	body := BindRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("参数错误: %s", err.Error())
		response.Fail(c)
		return
	}
	userID, _ := c.Get("user_id")

	// 将 userID 转换为 string 类型
	intUserID, ok := userID.(int)
	if !ok {
		// 转换失败，处理错误
		log.Println("Failed to convert user_id to string")
		// 进行错误处理逻辑...
		return
	}

	err := a.Service.BindUser(intUserID, body.Phone)

	if err != nil {
		fmt.Printf("请求错误: %s", err.Error())

		response.Fail(c)
		return
	}

	strUserID := strconv.Itoa(intUserID)

	// 在这里使用转换后的 string 类型的 userID 进行后续操作
	fmt.Println("User ID:", strUserID)

	data := gin.H{
		"message": "绑定成功",
	}
	response.OkWithData(data, c)
}

func (a Api) Ask(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		fmt.Printf("参数错误: %s", "code参数为空")
		response.FailWithMessage("code参数为空", c)
		return
	}

	println("获取参数code :", code)

	str, err := a.Service.Ask(code)

	println("!!")
	println(str)
	println("!!")

	if err != nil {
		fmt.Printf("请求错误: %s", err.Error())

		response.Fail(c)
		return
	}

	data := gin.H{
		"token": str,
	}
	response.OkWithData(data, c)
}

func (a Api) FindWxUser(c *gin.Context) {
	code := c.Query("openid")
	if code == "" {
		fmt.Printf("参数错误: %s", "openid参数为空")
		response.FailWithMessage("openid参数为空", c)
		return
	}

	println("获取参数openid :", code)

	user, err := a.Service.FindWxUser(code)

	if err != nil {
		fmt.Printf("请求错误: %s", err.Error())

		response.Fail(c)
		return
	}

	response.OkWithData(user, c)
}

func (a Api) Send(c *gin.Context) {
	// 使用 Query 方法获取查询参数值
	uidParam := c.Query("uid")

	// 将参数值转换为整数
	uid, err := strconv.Atoi(uidParam)
	if err != nil {
		// 处理转换错误，例如参数不是合法的整数
		c.JSON(400, gin.H{"error": "Invalid parameter"})
		return
	}

	err1 := a.Service.Send(uid)

	if err1 != nil {
		fmt.Printf("请求错误: %s", err1.Error())
		response.Fail(c)
		return
	}

	data := gin.H{
		"message": "发送成功",
	}
	response.OkWithData(data, c)
}

func (a Api) AskAbout(c *gin.Context) {
	body := AskRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Printf("参数错误: %s", err.Error())
		response.Fail(c)
		return
	}
	str, err := a.Service.Ask(body.Question)

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
