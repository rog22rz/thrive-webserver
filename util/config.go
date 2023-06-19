package util

import "github.com/spf13/viper"

type Config struct {
	APIToken            string `mapstructure:"API_TOKEN"`
	CreatorCollectionId string `mapstructure:"CREATOR_COLLECTION_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
