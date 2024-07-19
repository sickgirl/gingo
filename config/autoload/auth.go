package autoload

type User struct { //用户身份相关
	AuthKey string `mapstructure:"auth_key" json:"auth_key" yaml:"auth_key"` // 加密key
}
