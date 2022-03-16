package transfer

type UserState int

const (
	Idle UserState = iota
	PickingFirstService
)
