package http

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

type Response struct {
	Data       interface{}  `json:"data,omitempty"`
	HTTPCode   int          `json:"-"`
	Code       int          `json:"errno"`
	NewCode    int          `json:"code"`
	Message    string       `json:"errmsg"`
	NewMessage string       `json:"msg"`
	Type       ResponseType `json:"-"`
	BinaryBody []byte       `json:"-"`
	headers    map[string]string
}

type NewVersionResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
	Message string      `json:"msg"`
}

type ErrReason struct {
	OrderID int      `json:"id"`
	Reason  []string `json:"reason"`
}

type ResponseType int

const (
	_ ResponseType = iota
	ResponseTypeJson
	ResponseTypeImage
	ResponseTypeAttachment
	ResponseTypeRedirect
	ResponseTypeXML
)

func NewResponse() *Response {
	return &Response{
		HTTPCode: http.StatusOK,
		Code:     0,
		Message:  "",
		Type:     ResponseTypeJson,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func NewResponseWith(data interface{}) *Response {
	resp := NewResponse()
	resp.Data = data
	resp.NewCode = 200
	return resp
}

func NewResponseWithFile(data []byte, fileName string) *Response {
	return &Response{
		HTTPCode:   http.StatusOK,
		Code:       0,
		Message:    "",
		Type:       ResponseTypeAttachment,
		BinaryBody: data,
		headers: map[string]string{
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", fileName),
		},
	}
}

func (self *Response) Headers() map[string]string {
	return self.headers
}

func (self *Response) Write(w http.ResponseWriter) {
	for k, v := range self.headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(self.HTTPCode)

	switch self.Type {
	case ResponseTypeImage:
		{
			w.Write(self.BinaryBody)
		}
	case ResponseTypeAttachment:
		{
			w.Write(self.BinaryBody)
		}
	case ResponseTypeXML:
		{
			w.Write(self.BinaryBody)
		}
	case ResponseTypeRedirect:
	default:
		{
			body, err := jsoniter.Marshal(self)
			if err != nil {
			}
			w.Write(body)
		}
	}
}
