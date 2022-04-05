package mux

import "encoding/xml"

type Button struct {
	Text     string `xml:"text"`
	Metadata string `xml:"metadata"`
}

type URL struct {
	Link string `xml:"link"`
	Text string `xml:"text"`
}

type Either struct {
	Text   string `xml:"text"`
	URL    URL    `xml:"url"`
	Button Button `xml:"button"`
}

type Content struct {
	XMLName xml.Name `xml:"content"`
	Text    []string `xml:"text"`
	URL     []URL    `xml:"url"`
	Button  []Button `xml:"button"`
	Either  []Either `xml:"either"`
}
