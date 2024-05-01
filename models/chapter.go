package models

type Data struct {
	Chapter Chapter `json:"chapter"`
}

type Chapter struct {
	Number uint
	Name   string `json:"name"`
	Pages  []Page `json:"pages"`
}
