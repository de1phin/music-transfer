package mux

type Message struct {
	UserID  int64
	Content string
}

type UserState int

const (
	Idle UserState = iota
	ChoosingService
	ChoosingSrc
	ChoosingDst
)

type Transfer Service

type Song struct {
	Title   string
	Artists string
}

type Playlist struct {
	Title string
	Songs []Song
}
