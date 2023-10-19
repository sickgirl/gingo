package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	_ "time"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type Request struct {
	Method          string            `bson:"_"`
	Body            *bytes.Buffer     `bson:"-"`
	RemoteAddr      string            `bson:"client_ip,omitempty"`
	QueryParams     url.Values        `bson:"query_params,omitempty"`
	FormParams      url.Values        `bson:"form_params,omitempty"`
	URLParams       map[string]string `bson:"url_params,omitempty"`
	URL             *url.URL          `bson:"-"`
	MultipartForm   *multipart.Form   `bson:"-"`
	AccessToken     string            `bson:"access_token,omitempty"`
	Userhash        string            `bson:"userhash,omitempty"`
	DefaultLanguage string            `bson:"default_language,omitempty"`
	Languages       []string          `bson:"languages,omitempty"`
	UserAgent       string            `bson:"user_agent,omitempty"`
	OriginalRequest *http.Request     `bson:"-"`
}

func NewRequest(r *http.Request) (*Request, error) {
	var request Request
	request.Method = r.Method
	request.URL = r.URL
	request.Body = new(bytes.Buffer)
	request.Body.ReadFrom(r.Body)
	request.FormParams = r.Form
	request.RemoteAddr = r.Header.Get("X-Forwarded-For")
	request.MultipartForm = r.MultipartForm
	if request.RemoteAddr == "" {
		request.RemoteAddr = r.RemoteAddr
	}
	err := request.parseURL(r)
	if err != nil {
		return nil, err
	}
	err = request.parseHeader(r)
	if err != nil {
		return nil, err
	}
	request.OriginalRequest = r
	return &request, nil
}

func (r *Request) ShouldBind(obj interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, obj)
}

func (r *Request) GetStringFormParams(field string) (string, bool) {
	if r.FormParams[field] == nil {
		return "", false
	}
	return r.FormParams[field][0], true
}

// GetIntQueryParams 获取单个int params
func (r *Request) GetIntQueryParams(field string) (int, error) {
	if r.QueryParams[field] == nil {
		return 0, fmt.Errorf("invalid params %s", field)
	}
	i, err := strconv.Atoi(r.QueryParams[field][0])
	if err != nil {
		return 0, fmt.Errorf("invalid params %s", field)
	}
	return int(i), nil
}

// GetStringQueryParams 获取单个string params
func (r *Request) GetStringQueryParams(field string) (string, bool) {
	if r.QueryParams[field] == nil {
		return "", false
	}
	return r.QueryParams[field][0], true
}

func (r *Request) GetBoolQueryParams(field string) bool {
	if r.QueryParams[field] == nil {
		return false
	}
	return r.QueryParams[field][0] == "true"
}

func (r *Request) GetIntArrayQueryParams(field string) ([]int, error) {
	if r.QueryParams[field] == nil {
		return nil, fmt.Errorf("invalid params %s", field)
	}
	var values []int
	fields := r.QueryParams[field]
	for i := 0; i < len(fields); i++ {
		v, _ := strconv.Atoi(fields[i])
		values = append(values, v)
	}
	return values, nil
}

func (r *Request) parseURL(request *http.Request) error {
	r.URLParams = mux.Vars(request)
	r.QueryParams = r.URL.Query()
	return nil
}

func (r *Request) parseHeader(request *http.Request) error {
	r.AccessToken = request.Header.Get("x-access-token")
	r.Userhash = request.Header.Get("X-User-Userhash")
	language := request.Header.Get("Accept-Language")
	if language != "" {
		languages := strings.Split(language, ",")
		r.Languages = []string{}
		for i, v := range languages {
			if i == 0 {
				r.DefaultLanguage = strings.TrimSpace(v)
				if r.DefaultLanguage == "zh-CN" {
					r.DefaultLanguage = "zh"
				}
			}
			r.Languages = append(r.Languages, strings.TrimSpace(v))
		}
	}
	r.UserAgent = request.UserAgent()
	return nil
}
