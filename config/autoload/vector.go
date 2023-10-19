package autoload

type VECTOR struct { //向量数据库相关配置
	ApiKey     string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`          // apiKey
	Account    string `mapstructure:"account" json:"account" yaml:"account"`          // account
	Url        string `mapstructure:"url" json:"url" yaml:"url"`                      // url
	Db         string `mapstructure:"db" json:"db" yaml:"db"`                         // db
	Collection string `mapstructure:"collection" json:"collection" yaml:"collection"` // collection
}
