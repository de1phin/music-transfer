package transfer

type Interactor interface {
	GetMessageFrom() (int64, string)
	SendMessageTo(int64, string)
	ChooseFrom(int64, string, []string)
	SendURL(int64, string, string)
}
