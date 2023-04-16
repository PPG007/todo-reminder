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
	viper.Set("openai.sk", os.Getenv("OPENAI_SK"))
	viper.Set("email.password", os.Getenv("EMAIL_PASSWORD"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
