package transfer

type MusicService interface {
	Name() string
	URLName() string
	GetAuthURL(int64, string) string
}
