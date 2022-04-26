package spotify

type Config struct {
	RedirectURI  string `yaml:"redirectUri"`
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Scopes       string `yaml:"scopes"`
}
