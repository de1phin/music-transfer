package interactor

type Interactor interface {
	SendMessage(Message)
	GetMessage() Message
}

type Validator interface {
	Validate(*Message)
}

type InteractorSpec interface {
	Interactor
}

type interactorSpec struct {
	interactor Interactor
	validator  Validator
}

func NewInteractorSpec(interactor Interactor, validator Validator) *interactorSpec {
	return &interactorSpec{interactor: interactor, validator: validator}
}

func (is *interactorSpec) SendMessage(msg Message) {
	is.interactor.SendMessage(msg)
}

func (is *interactorSpec) GetMessage() (msg Message) {
	msg = is.interactor.GetMessage()
	is.validator.Validate(&msg)
	return msg
}
