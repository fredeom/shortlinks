package main

import (
	"html/template"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var id int = 0

type Link struct {
	Full  string
	Short string
	Id    int
}

func newLink(full, short string) Link {
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

func newData() Data {
	return Data{
		Links: []Link{
			newLink("https://www.google.com", "u/asdfaf"),
			newLink("https://yandex.ru", "u/zxvzxvc"),
		},
	}
}

func (d *Data) indexOf(id int) int {
	for i, link := range d.Links {
		if link.Id == id {
			return i
		}
	}
	return -1
}

func (d *Data) indexOfShortLink(shortlink string) int {
	for i, link := range d.Links {
		if link.Short == shortlink {
			return i
		}
	}
	return -1
}

func (d *Data) hasLink(fulllink string) bool {
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

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
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

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := newPage()

	e.Renderer = newTemplate()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.GET("/u/:param", func(c echo.Context) error {
		param := c.Param("param")
		idx := page.Data.indexOfShortLink("u/" + param)
		if idx == -1 {
			return c.String(400, "url not found")
		}

		return c.Redirect(301, page.Data.Links[idx].Full)
	})

	e.POST("/form", func(c echo.Context) error {
		fulllink := c.FormValue("fulllink")

		if page.Data.hasLink(fulllink) {
			formData := newFormData()
			formData.Values["fulllink"] = fulllink
			formData.Errors["fulllink"] = "Fulllink already exists"

			return c.Render(422, "form", formData)
		}

		for {
			rs := RandStringBytes(6)
			if !page.Data.hasLink(rs) {
				link := newLink(fulllink, "u"+"/"+rs)
				page.Data.Links = append(page.Data.Links, link)

				c.Render(200, "form", newFormData())
				return c.Render(200, "oob-link", link)
			}
		}
	})

	e.DELETE("/links/:id", func(c echo.Context) error {
		time.Sleep(1 * time.Second)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(400, "Invalid id")
		}

		index := page.Data.indexOf(id)
		if index == -1 {
			return c.String(404, "Link not found")
		}

		page.Data.Links = append(page.Data.Links[:index], page.Data.Links[index+1:]...)

		return c.NoContent(200)

	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
