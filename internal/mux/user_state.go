package mux

type UserState int

const (
	Idle UserState = iota
	ChoosingService
)
