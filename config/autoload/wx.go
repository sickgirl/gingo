package autoload

type WX struct { //微信相关配置
	AppID     string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`             // AppID
	AppSecret string `mapstructure:"app_secret" json:"app_secret" yaml:"app_secret"` // AppSecret
	Url       string `mapstructure:"url" json:"url" yaml:"url"`                      // Url
}
