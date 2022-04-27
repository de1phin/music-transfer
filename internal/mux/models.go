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

type TransferState struct {
	sourceServiceName       string
	sourceServiceAuthorized bool
	destinationServiceName  string
	activeInteractorName    string
	interactorUserID        int64
}

type UserState int

const (
	Idle UserState = iota
	ChooseSource
	AuthorizeSource
	ChooseDestination
	AuthorizeDestination
	Transfer
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
