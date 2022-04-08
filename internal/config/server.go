package config

import "github.com/spf13/viper"

func (*config) GetServerURL() string {
	return viper.GetString("server.URL")
}

func (*config) GetServerHostname() string {
	return viper.GetString("server.hostname")
}
