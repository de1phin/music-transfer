package config

import (
	"github.com/spf13/viper"
)

func (*config) GetYouTubeApiKEY() string {
	return viper.GetString("youtube.APIKey")
}

func (*config) GetYouTubeClientID() string {
	return viper.GetString("youtube.clientID")
}

func (*config) GetYouTubeClientSecret() string {
	return viper.GetString("youtube.clientSecret")
}

func (*config) GetYouTubeScope() string {
	return viper.GetString("youtube.scope")
}

func (*config) GetYouTubeEndpoint() string {
	return viper.GetString("youtube.endpoint")
}
