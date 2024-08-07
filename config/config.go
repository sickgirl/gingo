package config

import (
	"github.com/robfig/cron/v3"
	"github.com/songcser/gingo/config/autoload"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Configuration struct {
	Domain string          `mapstructure:"domain" json:"domain" yaml:"domain"`
	DbType string          `mapstructure:"dbType" json:"dbType" yaml:"dbType"`
	Admin  autoload.Admin  `mapstructure:"admin" json:"admin" yaml:"admin"`
	Mysql  autoload.Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap    autoload.Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	JWT    autoload.JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	BC     autoload.BC     `mapstructure:"bc" json:"bc" yaml:"bc"`
	WX     autoload.WX     `mapstructure:"wx" json:"wx" yaml:"wx"`
	VECTOR autoload.VECTOR `mapstructure:"vector" json:"vector" yaml:"vector"`
	User   autoload.User   `mapstructure:"user" json:"user" yaml:"user"`
}

var (
	GVA_CONFIG Configuration
	GVA_DB     *gorm.DB
	GVA_LOG    *zap.Logger
	GVA_VP     *viper.Viper
	GVA_JOB    *cron.Cron
)
