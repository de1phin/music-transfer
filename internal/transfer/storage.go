package transfer

type Storage interface {
	AddService(string, interface{})
	PutServiceData(string, int64, interface{})
	GetServiceData(string, int64) interface{}
}
