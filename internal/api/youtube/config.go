package youtube

type Config struct {
	APIKey       string `yaml:"apiKey"`
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Scopes       string `yaml:"scopes"`
	RedirectURI  string `yaml:"redirectUri"`
}
