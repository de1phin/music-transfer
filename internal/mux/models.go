package mux

type Message struct {
	UserID    int64
	UserState UserState
	Content   MessageContent
}

type MessageContent struct {
	Text    string
	URLs    []URL
	Buttons []string
}

type Transfer Service

type UserState int

const (
	Idle UserState = iota
	ChoosingSrc
	ChoosingDst
)

type URL struct {
	Link string
	Text string
}

type Song struct {
	Title   string
	Artists string
}

type Playlist struct {
	Title string
	Songs []Song
}
