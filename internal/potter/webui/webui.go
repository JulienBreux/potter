package webui

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/julienbreux/potter/pkg/color"
	"github.com/julienbreux/potter/webui"
)

func New(version string) {
	bgcolor := color.RandomColor()
	engine := html.NewFileSystem(http.FS(webui.Views), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"version": version,
			"bgColor": bgcolor,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
