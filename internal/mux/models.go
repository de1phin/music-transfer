package mux

import "encoding/xml"

type Message struct {
	UserID    int64
	UserState UserState
	Content   string
}

type UserState int

const (
	Idle UserState = iota
	ChoosingService
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

type Content struct {
	XMLName xml.Name `xml:"content"`
	Text    []string `xml:"text"`
	URL     []URL    `xml:"url"`
	Button  []string `xml:"button"`
	Either  []Either `xml:"either"`
}

type Transfer Service

type Song struct {
	Title   string
	Artists string
}

type Playlist struct {
	Title string
	Songs []Song
}
