package transfer

type Interactor interface {
	GetMessageFrom() (string, int64)
	SendMessageTo(string, int64)
}
