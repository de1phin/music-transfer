package config

import "github.com/spf13/viper"

func (*config) GetSpotifyClientID() string {
	return viper.GetString("spotify.clientID")
}

func (*config) GetSpotifyScope() string {
	return viper.GetString("spotify.scope")
}

func (*config) GetSpotifyClientSecret() string {
	return viper.GetString("spotify.clientSecret")
}

func (*config) GetSpotifyEndpoint() string {
	return viper.GetString("spotify.endpoint")
}
