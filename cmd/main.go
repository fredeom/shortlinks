package main

import (
	"os"

	"github.com/fredeom/shortlinks/db"
	"github.com/fredeom/shortlinks/handlers"
	"github.com/fredeom/shortlinks/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const dbName = "data.db"

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/images", "images")
	e.Static("/css", "css")

	lStore, err := db.NewLinkStore(dbName)
	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	formData := services.NewFormData()

	ls := services.NewServicesLink(formData, lStore)

	h := handlers.New(ls)

	handlers.SetupRoutes(e, h)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
