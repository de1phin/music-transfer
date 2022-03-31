package config

import "github.com/spf13/viper"

func (*config) GetCallbackURL() string {
	return viper.GetString("callbackURL")
}

func (*config) GetServerURL() string {
	return viper.GetString("serverURL")
}
