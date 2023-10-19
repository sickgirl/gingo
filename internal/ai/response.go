package ai

type Response struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Type        string `json:"type"`
}
type BQResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  Data   `json:"data"`
	Usage Usage  `json:"usage"`
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
