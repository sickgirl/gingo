package autoload

type VECTOR struct { //向量数据库相关配置
	ApiKey    string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`          // apiKey
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"` // secretKey
	Url       string `mapstructure:"url" json:"url" yaml:"url"`                      // url
}
