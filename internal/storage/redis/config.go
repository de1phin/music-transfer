package redis

type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}
