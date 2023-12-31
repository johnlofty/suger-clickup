package settings

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

type config struct {
	ClickupConfig clickupConfig `mapstructure:"clickup"`
	DBConfig      dbConfig      `mapstructure:"database"`
	JWT           jwtConfig     `mapstructure:"jwt"`
}

var conf *config

type dbConfig struct {
	Host     string
	DBName   string
	Port     int
	User     string
	Password string
	Sslmode  string
}

type clickupConfig struct {
	Secret string
	ListId string
}

type jwtConfig struct {
	Secret string
}

func (c dbConfig) URI() string {
	return fmt.Sprintf("postgress://%s:%s@%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.DBName)
}

func (c dbConfig) Format() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		c.User, c.Password, c.DBName, c.Host)
}

func Setup() {
	SetupLog()
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Println("read config error:", err)
		panic(err)
	}
	log.Debugf("reading config:%#+v\n", conf)
}

func Get() *config {
	return conf
}

func SetupLog() {
	log.SetLevel(log.LevelDebug)
}
