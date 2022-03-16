package transfer

type Storage interface {
	GetUserState(int64) UserState
	PutUserState(int64, UserState)
	AddService(string)
	PutServiceData(int64, string, interface{})
	GetServiceData(int64, string) interface{}
}
