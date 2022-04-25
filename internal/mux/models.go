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

type Either struct {
	Text   string `xml:"text"`
	URL    URL    `xml:"url"`
	Button string `xml:"button"`
}

type URL struct {
	Link string `xml:"link"`
	Text string `xml:"text"`
}

type Song struct {
	Title   string
	Artists string
}

type Playlist struct {
	Title string
	Songs []Song
}
