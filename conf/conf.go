package conf

import (
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("application")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	viper.Set("minio.ak", os.Getenv("MINIO_AK"))
	viper.Set("minio.sk", os.Getenv("MINIO_SK"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
