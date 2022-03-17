package transfer

type Config interface {
	GetCallbackURL() string
}
