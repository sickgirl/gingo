package ai

type Request struct {
	Name  string `json:"name" form:"name"`
	Level string `json:"level" form:"level"`
}
