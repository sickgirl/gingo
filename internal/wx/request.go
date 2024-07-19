package wx

type Request struct {
	Name  string `json:"name" form:"name"`
	Level string `json:"level" form:"level"`
}

type AskRequest struct {
	Question string `json:"question" form:"question"`
}

type VectorRequest struct {
	Database   string `json:"database"`
	Collection string `json:"collection"`
	Search     Search `json:"search"`
}

type BindRequest struct {
	Phone string `json:"phone"`
}

//SubscribeMessageRequest 订阅消息发送请求的结构体
type SubscribeMessageRequest struct {
	ToUser           string                 `json:"touser"`
	TemplateID       string                 `json:"template_id"`
	Page             string                 `json:"page"`
	Data             map[string]interface{} `json:"data"`
	MiniprogramState string                 `json:"miniprogram_state"`
}

type Params struct {
	Ef int `json:"ef"`
}

type Search struct {
	EmbeddingItems []string `json:"embeddingItems"`
	Limit          int      `json:"limit"`
	Params         Params   `json:"params"`
	RetrieveVector bool     `json:"retrieveVector"`
	Filter         string   `json:"filter"`
	OutputFields   []string `json:"outputFields"`
}

type BCRequest struct {
	Model      string      `json:"model"`
	Messages   []MessagesB `json:"messages"`
	Parameters ParametersB `json:"parameters"`
}

type MessagesB struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ParametersB struct {
	Temperature float64 `json:"temperature"`
	TopK        int     `json:"top_k"`
}
