package conf

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("application")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
