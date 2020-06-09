package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

var(
	config = "config"
	configPath = "."
)

func InitConfiguration() {
	viper.SetConfigName(config)
	viper.AddConfigPath(configPath)
	if err :=viper.ReadInConfig(); err!= nil{
		log.Fatalln("Exception while init config", err)
	}
}

func InitLogConfig() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}
