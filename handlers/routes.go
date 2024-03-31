package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo, h *LinkHandler) {
	app.GET("/", h.HandleIndex)
	app.GET("/u/:param", h.HandleShortUrl)
	app.POST("/form", h.HandleForm)
	app.DELETE("/links/:id", h.HandleDelete)
	app.GET("/visitors/:id", h.HandleVisitors)
}
