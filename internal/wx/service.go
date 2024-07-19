package wx

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/http"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/utils"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"time"
)

//Parameters 定义包含参数的结构体
type Parameters struct {
	Temperature float64 `json:"temperature"`
	TopK        int     `json:"top_k"`
}

type Service struct {
	service.Service[TWxUsers]
}

func NewService(a TWxUsers) Service {
	return Service{service.NewBaseService[TWxUsers](a)}
}

func (s Service) MakeMapper(c *gin.Context) model.Mapper[TWxUsers] {
	var r Request
	err := c.ShouldBindQuery(&r)
	utils.CheckError(err)
	w := model.NewWrapper()
	w.Like("name", r.Name)
	w.Eq("level", r.Level)
	m := model.NewMapper[TWxUsers](TWxUsers{}, w)
	return m
}

//Embedding  将问题向量化 并从数据库中查询出匹配数据
func (s Service) Embedding(question string) (string, error) {

	messAge := ""

	VConf := config.GVA_CONFIG.VECTOR
	apiKey := VConf.ApiKey
	account := VConf.Account
	db := VConf.Db
	collection := VConf.Collection
	url := VConf.Url

	request := VectorRequest{
		Database:   db,
		Collection: collection,
		Search:     Search{},
	}
	req, _ := json.Marshal(request)

	header := http.NewHeader()
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+fmt.Sprintf("account=%s&api_key=%s", account, apiKey))

	resp, err := http.Call("POST", url, header, req)
	if err != nil {
		return messAge, err
	}

	data, err2 := ioutil.ReadAll(resp.Data.(io.Reader))
	if err2 != nil {
		return messAge, err2
	}
	var res VResponse
	if err := jsoniter.Unmarshal(data, &res); err != nil {
		return messAge, err
	}

	for _, v := range res.Documents {
		messAge += v.Text
	}
	messAge = messAge + " " + question
	return messAge, nil
}

//CheckFromWx  传入匹配完的数据文本 并返回答案
func (s Service) CheckFromWx(code string) (string, error) {
	// API Key 和 Secret Key  及url
	WXConf := config.GVA_CONFIG.WX
	appId := WXConf.AppID
	appSecret := WXConf.AppSecret
	url := WXConf.Url

	queryUrl := fmt.Sprintf(url+"/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, appSecret, code)

	// 计算 X-BC-Timestamp（UTC时间戳）
	//timestamp := time.Now().Unix()

	header := http.NewHeader()
	header.Set("Content-Type", "application/json")

	fmt.Println("请求接口地址:", queryUrl)

	resp, err := http.Call("GET", queryUrl, header, []byte{})
	if err != nil {
		fmt.Println("请求接口错误!:")
		fmt.Println("请求接口地址:", queryUrl)
		fmt.Println("请求接口错误:", err.Error())

		return "", err
	}

	body, err := ioutil.ReadAll(resp.Data.(io.Reader))
	if err != nil {
		fmt.Println("请求接口读取内容错误!:")
		fmt.Println("请求接口读取内容错误!:Error:", err.Error())
		return "", err
	}
	println("请求微信接口获取当前登录用户信息  数据如下:")

	fmt.Printf("%+v", string(body))
	println(body)

	var res ResponseData
	if err := jsoniter.Unmarshal(body, &res); err != nil {
		fmt.Println("请求接口格式化内容错误!:")
		fmt.Println("请求接口格式化内容错误!:Error:", err.Error())
		return "", err
	}

	if res.ErrCode != SuccessErrCode {
		fmt.Println("微信返回数据错误!:")
		fmt.Println("微信返回数据错误:", res.ErrMsg)

		fmt.Printf("%v", res)

		return "", errors.New(res.ErrMsg)
	}

	var Nid int64
	var err1 error

	user, err := new(TWxUsers).FindOneByOpenid(res.OpenID)
	if user.ID > 0 {
		println("数据存在!")
		Nid = user.ID
	} else {
		if err == gorm.ErrRecordNotFound {
			createdAt := time.Now()

			newUser := TWxUsers{
				BaseModel: model.BaseModel{
					CreatedAt: utils.JsonTime(createdAt),
				},
				Openid:  res.OpenID,
				Appid:   1,
				Unionid: res.UnionID,
			}

			Nid, err1 = newUser.Create()
			if err1 != nil {
				config.GVA_LOG.Error("插入数据异常:" + err.Error())
			}
		}
	}

	return GenerateToken(int(Nid), "当前用户信息xxx")
}

func getAccessToken() (string, error) {
	// 在这里实现获取访问令牌的逻辑
	// 可以使用微信提供的凭证获取接口，或者使用第三方库等方法获取
	// ...
	WXConf := config.GVA_CONFIG.WX
	appID := WXConf.AppID
	appSecret := WXConf.AppSecret
	url := WXConf.Url
	queryUrl := fmt.Sprintf(url+"/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appID, appSecret)

	header := http.NewHeader()
	header.Set("Content-Type", "application/json")

	fmt.Println("请求接口地址:", queryUrl)

	resp, err := http.Call("GET", queryUrl, header, []byte{})
	if err != nil {
		fmt.Println("请求接口错误!:")
		fmt.Println("请求接口地址:", queryUrl)
		fmt.Println("请求接口错误:", err.Error())

		return "", err
	}

	body, err := ioutil.ReadAll(resp.Data.(io.Reader))
	if err != nil {
		fmt.Println("请求接口读取内容错误!:")
		fmt.Println("请求接口读取内容错误!:Error:", err.Error())
		return "", err
	}

	var res ResponseTokenData
	if err := jsoniter.Unmarshal(body, &res); err != nil {
		fmt.Println("请求接口格式化内容错误!:")
		fmt.Println("请求接口格式化内容错误!:Error:", err.Error())
		return "", err
	}

	if res.ErrCode != SuccessErrCode {
		fmt.Println("微信返回数据错误!:")
		fmt.Println("微信返回数据错误:", res.ErrMsg)

		fmt.Printf("%v", res)

		return "", errors.New(res.ErrMsg)
	}

	return res.AccessToken, nil

}

//SendMsgByWx  通过微信发送订阅消息
func (s Service) SendMsgByWx(openid string) error {
	// API Key 和 Secret Key  及url
	WXConf := config.GVA_CONFIG.WX
	//appId := WXConf.AppID
	//appSecret := WXConf.AppSecret
	url := WXConf.Url

	// 构建请求体
	requestBody := SubscribeMessageRequest{
		ToUser:           openid,
		TemplateID:       "zuAdu8G9DziAsxcPjMgkkWgVvOCzkEmygrhfqPjezgo",
		Page:             "pages/scheduling1/scheduling1", // 替换为跳转的页面路径
		MiniprogramState: "trial",
		Data: map[string]interface{}{
			"thing1": map[string]string{
				"value": "高山",
			},
			"thing4": map[string]string{
				"value": "周一",
			},
		},
	}

	token, err := getAccessToken()

	if err != nil {
		return err
	}

	// 构建请求的 URL
	queryUrl := fmt.Sprintf(url+"/cgi-bin/message/subscribe/send?access_token=%s", token)

	// 计算 X-BC-Timestamp（UTC时间戳）
	//timestamp := time.Now().Unix()

	header := http.NewHeader()
	header.Set("Content-Type", "application/json")

	fmt.Println("请求接口地址:", queryUrl)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Call("POST", queryUrl, header, jsonData)
	if err != nil {
		fmt.Println("请求接口错误!:")
		fmt.Println("请求接口地址:", queryUrl)
		fmt.Println("请求接口错误:", err.Error())

		return err
	}

	body, err := ioutil.ReadAll(resp.Data.(io.Reader))
	if err != nil {
		fmt.Println("请求接口读取内容错误!:")
		fmt.Println("请求接口读取内容错误!:Error:", err.Error())
		return err
	}

	var res SendResponseData
	if err := jsoniter.Unmarshal(body, &res); err != nil {
		fmt.Println("请求接口格式化内容错误!:")
		fmt.Println("请求接口格式化内容错误!:Error:", err.Error())
		return err
	}

	if res.ErrCode != SuccessErrCode {
		fmt.Println("微信返回数据错误!:")
		fmt.Println("微信返回数据错误:", res.ErrMsg)

		fmt.Printf("%v", res)

		return errors.New(res.ErrMsg)
	}

	return nil
}

func (s Service) Test(question string) (string, error) {
	return "测试返回", nil
}

//BindUser 通过员工手机号匹配 绑定微信用户 及 员工表
func (s Service) BindUser(wxUid int, phone string) error {
	//todo 根据 phone 查找员工
	// 查到员工  并且关联关系不存在时候 进行绑定  把wxUuid 写入绑定关系里
	// 操作成功 返回

	user, err := new(TUsers).FindOneByPhone(phone)
	if err != nil {
		return err
	}

	if user.WxId != 0 {
		return errors.New("该用户已经绑定")
	}

	err1 := new(TUsers).BindUser(phone, wxUid)
	if err1 != nil {
		return err1
	}
	return nil
}

func (s Service) Ask(code string) (string, error) {
	return s.CheckFromWx(code)
}

func (s Service) FindWxUser(openid string) (TWxUsers, error) {

	user, err := new(TWxUsers).FindOneByOpenid(openid)
	if user.ID > 0 {
		println("数据存在!")
	}

	openid2 := "test:openid"
	unionid := "test:unionid"
	createdAt := time.Now()

	newUser := TWxUsers{
		BaseModel: model.BaseModel{
			CreatedAt: utils.JsonTime(createdAt),
		},
		Openid:  openid2,
		Appid:   0,
		Unionid: unionid,
	}

	if err == gorm.ErrRecordNotFound {
		id, err1 := newUser.Create()
		if err1 != nil || id <= 0 {
			config.GVA_LOG.Error("插入数据异常:" + err.Error())
		}
	}

	return new(TWxUsers).FindOneByOpenid(openid)
}

func (s Service) Send(uid int) error {
	//todo 通过数据库  用uid 获取对应的 openid

	openid := "osb3V4nsDgSfyTt2_r5V0U8FGbos"
	return s.SendMsgByWx(openid)
}

//func (s Service) AskAbout(question string) (string, error) {
//
//	question, err := s.Embedding(question)
//	if err != nil {
//		return "", err
//	}
//
//	return s.CheckFromBC(question)
//}

// 计算签名
func calculateSignature(secretKey, requestBody string, timestamp int64) string {
	dataToSign := secretKey + requestBody + fmt.Sprint(timestamp)
	hasher := md5.New()
	hasher.Write([]byte(dataToSign))
	return hex.EncodeToString(hasher.Sum(nil))
}
