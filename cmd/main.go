package main

import (
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/fredeom/shortlinks/services"
	"github.com/fredeom/shortlinks/views"
)

func newOutput(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := services.NewPage()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return newOutput(c, views.Index(page))
	})

	e.GET("/u/:param", func(c echo.Context) error {
		param := c.Param("param")
		idx := page.Data.IndexOfShortLink("u/" + param)
		if idx == -1 {
			return c.String(400, "url not found")
		}

		return c.Redirect(301, page.Data.Links[idx].Full)
	})

	e.POST("/form", func(c echo.Context) error {
		fulllink := c.FormValue("fulllink")

		if page.Data.HasLink(fulllink) {
			formData := services.NewFormData()
			formData.Values["fulllink"] = fulllink
			formData.Errors["fulllink"] = "Fulllink already exists"

			c.Response().Status = 422
			return newOutput(c, views.Form(formData))
		}

		for {
			rs := services.RandStringBytes(6)
			if !page.Data.HasLink(rs) {
				link := services.NewLink(fulllink, "u"+"/"+rs)
				page.Data.Links = append(page.Data.Links, link)

				newOutput(c, views.Form(services.NewFormData()))
				return newOutput(c, views.OobLink(link))
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

		index := page.Data.IndexOf(id)
		if index == -1 {
			return c.String(404, "Link not found")
		}

		page.Data.Links = append(page.Data.Links[:index], page.Data.Links[index+1:]...)

		return c.NoContent(200)

	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
