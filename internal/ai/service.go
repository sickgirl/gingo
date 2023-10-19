package ai

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/songcser/gingo/config"
	"github.com/songcser/gingo/pkg/http"
	"github.com/songcser/gingo/pkg/model"
	"github.com/songcser/gingo/pkg/service"
	"github.com/songcser/gingo/utils"
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
	service.Service[App]
}

func NewService(a App) Service {
	return Service{service.NewBaseService[App](a)}
}

func (s Service) MakeMapper(c *gin.Context) model.Mapper[App] {
	var r Request
	err := c.ShouldBindQuery(&r)
	utils.CheckError(err)
	w := model.NewWrapper()
	w.Like("name", r.Name)
	w.Eq("level", r.Level)
	m := model.NewMapper[App](App{}, w)
	return m
}

func (s Service) MakeResponse(val model.Model) any {
	a := val.(App)
	res := Response{
		Name:        a.Name,
		Description: fmt.Sprintf("名称：%s, 等级: %s, 类型: %s", a.Name, a.Level, a.Type),
		Level:       a.Level,
		Type:        a.Type,
	}
	return res
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

//CheckFromBC  传入匹配完的数据文本 并返回答案
func (s Service) CheckFromBC(question string) (string, error) {
	// API Key 和 Secret Key  及url
	BCconf := config.GVA_CONFIG.BC
	apiKey := BCconf.ApiKey
	secretKey := BCconf.SecretKey
	url := BCconf.Url

	// 计算 X-BC-Timestamp（UTC时间戳）
	timestamp := time.Now().Unix()

	var mySlice []MessagesB
	element1 := MessagesB{
		Role:    "user",
		Content: question,
	}

	mySlice = append(mySlice, element1)

	messAge := ""
	request := BCRequest{
		Model:    "Baichuan2-53B",
		Messages: mySlice,
		Parameters: ParametersB{
			Temperature: 0.3,
			TopK:        10,
		},
	}
	req, _ := json.Marshal(request)

	// 计算签名（X-BC-Signature）
	signature := calculateSignature(secretKey, string(req), timestamp)

	header := http.NewHeader()
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "Bearer "+apiKey)
	header.Set("X-BC-Timestamp", fmt.Sprint(timestamp))
	header.Set("X-BC-Signature", signature)
	header.Set("X-BC-Sign-Algo", "MD5")

	resp, err := http.Call("POST", url, header, req)
	if err != nil {
		return messAge, err
	}
	if !(resp.HTTPCode >= 200 && resp.HTTPCode <= 299) {
		return messAge, fmt.Errorf("Request error %v", resp.HTTPCode)
	}

	data, err := ioutil.ReadAll(resp.Data.(io.Reader))
	if err != nil {
		return messAge, err
	}
	var res BQResponse
	if err := jsoniter.Unmarshal(data, &res); err != nil {
		return messAge, err
	}

	messAge = res.Data.Messages[0].Content

	return messAge, nil
}

func (s Service) Test(question string) (string, error) {
	return "测试返回", nil
}

func (s Service) Ask(question string) (string, error) {
	return s.CheckFromBC(question)
}

func (s Service) AskAbout(question string) (string, error) {

	question, err := s.Embedding(question)
	if err != nil {
		return "", err
	}

	return s.CheckFromBC(question)
}

// 计算签名
func calculateSignature(secretKey, requestBody string, timestamp int64) string {
	dataToSign := secretKey + requestBody + fmt.Sprint(timestamp)
	hasher := md5.New()
	hasher.Write([]byte(dataToSign))
	return hex.EncodeToString(hasher.Sum(nil))
}
