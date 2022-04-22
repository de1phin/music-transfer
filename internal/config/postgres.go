package config

import "github.com/spf13/viper"

func (*config) GetPosgresDataSourceName() string {
	return viper.GetString("postgres.dataSourceName")
}
