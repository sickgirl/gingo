package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

type Header struct {
	http.Header
}

func NewHeader() *Header {
	return &Header{make(http.Header)}
}

func (self *Header) Set(key, value string) *Header {
	self.Header.Set(key, value)
	return self
}

func (self *Header) SetAccessToken(accessToken string) *Header {
	self.Header.Set("X-User-AccessToken", accessToken)
	return self
}

type URL struct {
	*url.URL
	query url.Values
}

func NewURL(rawurl string) *URL {
	var u URL
	var err error
	u.URL, err = url.Parse(rawurl)
	if err != nil {
	} else {
		u.query = u.URL.Query()
	}
	return &u
}

func (self *URL) Set(key, value string) *URL {
	self.query.Set(key, value)
	return self
}

func (self *URL) String() string {
	self.URL.RawQuery = self.query.Encode()
	return self.URL.String()
}

func GetHost(hostPort string) string {
	if index := strings.Index(hostPort, ":"); index != -1 {
		return hostPort[:index]
	} else {
		return hostPort
	}
}

func JsonToReader(v interface{}) io.Reader {
	data, err := jsoniter.Marshal(v)
	if err != nil {
		return nil
	}
	return bytes.NewReader(data)
}

func PostJson(url string, requestJson, responseJson interface{}) (*http.Response, error) {
	resp, err := http.Post(url, "application/json", JsonToReader(requestJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := jsoniter.NewDecoder(resp.Body)
	err = decoder.Decode(responseJson)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func GetJson(url string, responseJson interface{}) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := jsoniter.NewDecoder(resp.Body)
	err = decoder.Decode(responseJson)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func InternalResponseData(resp *http.Response) (*json.RawMessage, error) {
	var jsonRaw map[string]*json.RawMessage
	bytes, _ := ioutil.ReadAll(resp.Body.(io.Reader))
	if err := jsoniter.Unmarshal(bytes, &jsonRaw); err != nil {
		return nil, err
	} else {
		return jsonRaw["data"], nil
	}
}

func Call(method, url string, header *Header, requestBody []byte) (*Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = header.Header
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()

	if resp.StatusCode != 200 {
	}

	var response Response
	response.Data = resp.Body

	response.HTTPCode = resp.StatusCode
	return &response, nil
}
