package config

import (
	"github.com/spf13/viper"
	"log"
)

type LogConfig struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`
	Format       string `mapstructure:"format" json:"format" yaml:"format"`
	LogDir       string `mapstructure:"log_dir" json:"log_dir" yaml:"log_dir"`
	ShowLine     bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	LogInConsole bool   `mapstructure:"log_in_console" json:"log_in_console" yaml:"log_in_console"`
	MaxSize      int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge       int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	Compress     bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type Document struct {
	SecretId  string `mapstructure:"secret_id" json:"secret_id" yaml:"secret_id"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	AppId     string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`
	AppKey    string `mapstructure:"app_key" json:"app_key" yaml:"app_key"`
}

type AppConfig struct {
	App struct {
		Port int `mapstructure:"port"`
	}

	UserFile   string `mapstructure:"user_file"`
	CookieFile string `mapstructure:"cookie_file"`

	Log LogConfig `mapstructure:"log"`

	Document Document `mapstructure:"document"`
}

var GlobalConfig AppConfig

func Init() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	log.Println("配置加载成功")
}
