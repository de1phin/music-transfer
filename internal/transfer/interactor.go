package transfer

type Interactor interface {
	GetMessageFrom() (int64, string)
	SendMessageTo(int64, string)
}
