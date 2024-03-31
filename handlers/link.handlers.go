package handlers

import (
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/fredeom/shortlinks/services"
	"github.com/fredeom/shortlinks/views"
	"github.com/labstack/echo/v4"
)

type LinkService interface {
	GetFormData() (services.FormData, error)
	GetData() (services.Data, error)
	HasLink(s string) int
	GetLink(id int) services.Link
	NewLink(full, short string) services.Link
	DeleteLink(services.Link)
	UpdateHits(id int)
	AddVisitor(id int, userAgent string)
	GetVisitors(id int) services.Visitors
}

type LinkHandler struct {
	LinkService LinkService
}

func New(ls LinkService) *LinkHandler {
	return &LinkHandler{
		LinkService: ls,
	}
}

func (lh *LinkHandler) HandleIndex(c echo.Context) error {
	formData, _ := lh.LinkService.GetFormData()
	data, _ := lh.LinkService.GetData()
	return lh.View(c, views.Index(formData, data))
}

func (lh *LinkHandler) HandleShortUrl(c echo.Context) error {
	param := c.Param("param")
	idx := lh.LinkService.HasLink("u/" + param)
	if idx == -1 {
		return c.String(400, "url not found")
	}
	lh.LinkService.UpdateHits(idx)
	lh.LinkService.AddVisitor(idx, c.Request().UserAgent())
	return c.Redirect(302, lh.LinkService.GetLink(idx).Full)
}

func (lh *LinkHandler) HandleForm(c echo.Context) error {
	fulllink := c.FormValue("fulllink")

	if lh.LinkService.HasLink(fulllink) >= 0 {
		formData := services.NewFormData()
		formData.Values["fulllink"] = fulllink
		formData.Errors["fulllink"] = "Fulllink already exists"

		c.Response().Status = 422
		return lh.View(c, views.Form(formData))
	}
	for {
		rs := services.RandStringBytes(6)
		if lh.LinkService.HasLink("u"+"/"+rs) == -1 {
			link := lh.LinkService.NewLink(fulllink, "u"+"/"+rs)
			lh.View(c, views.Form(services.NewFormData()))
			return lh.View(c, views.OobLink(link))
		}
	}
}

func (lh *LinkHandler) HandleDelete(c echo.Context) error {
	time.Sleep(1 * time.Second)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "Invalid id")
	}

	link := lh.LinkService.GetLink(id)
	if link.ID == -1 {
		return c.String(404, "Link not found")
	}

	lh.LinkService.DeleteLink(link)

	return c.NoContent(200)
}

func (lh *LinkHandler) HandleVisitors(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "Invalid id")
	}

	return lh.View(c, views.Visitors(lh.LinkService.GetVisitors(id)))
}

func (lh *LinkHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
