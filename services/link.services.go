package services

import (
	"math/rand"
)

var id int = 0

type Link struct {
	Full  string
	Short string
	Id    int
}

func NewLink(full, short string) Link {
	id++
	return Link{
		Full:  full,
		Short: short,
		Id:    id,
	}
}

type Links = []Link

type Data struct {
	Links Links
}

func NewData() Data {
	return Data{
		Links: []Link{
			NewLink("https://www.google.com", "u/asdfaf"),
			NewLink("https://yandex.ru", "u/zxvzxvc"),
		},
	}
}

func (d *Data) IndexOf(id int) int {
	for i, link := range d.Links {
		if link.Id == id {
			return i
		}
	}
	return -1
}

func (d *Data) IndexOfShortLink(shortlink string) int {
	for i, link := range d.Links {
		if link.Short == shortlink {
			return i
		}
	}
	return -1
}

func (d *Data) HasLink(fulllink string) bool {
	for _, link := range d.Links {
		if link.Full == fulllink || link.Short == fulllink {
			return true
		}
	}
	return false
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func NewPage() Page {
	return Page{
		Data: NewData(),
		Form: NewFormData(),
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
