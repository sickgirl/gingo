package wx

const SuccessErrCode = 0

type Response struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Type        string `json:"type"`
}
type ResponseData struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type ResponseTokenData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
}

type SendResponseData struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type Messages struct {
	Role         string `json:"role"`
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
}

type Data struct {
	Messages []Messages `json:"messages"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	AnswerTokens int `json:"answer_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type VResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Documents []Documents `json:"documents"`
}

type Documents struct {
	ID       string  `json:"id"`
	Score    float64 `json:"score"`
	BookName string  `json:"bookName"`
	Author   string  `json:"author"`
	Text     string  `json:"text"`
}
