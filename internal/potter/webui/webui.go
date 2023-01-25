package webui

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/julienbreux/potter/pkg/color"
	"github.com/julienbreux/potter/pkg/emoji"
	"github.com/julienbreux/potter/pkg/namesgen"
)

const (
	// TODO: Move to configuration
	port = 8080

	dateFormat = "02-01-2006 15:04:05"

	directory = "./internal/potter/webui/views"
	extension = ".go.html"
)

// WebUI represents a Web UI server
type Webui struct {
	app *fiber.App
}

// New creates a new Web UI server
func New(version string) Webui {
	now := time.Now()

	engine := html.New(directory, extension)
	// TODO: Add to configuration for dev
	engine.Reload(true)
	engine.AddFunc("currentTime", func() string {
		return time.Now().Format(dateFormat)
	})

	config := fiber.Config{
		DisableStartupMessage: true,

		Views: engine,
	}
	app := fiber.New(config)

	// TODO: Move to config
	color := color.Rand()
	bgColor := color.ToHex()
	fgColor := color.Invert().ToHex()
	name := namesgen.GetRandom()
	emoji := emoji.GetRandom()

	app.Static("assets", "./assets")

	app.Get("/", func(c *fiber.Ctx) error {
		ua := string(c.Context().UserAgent())
		vars := fiber.Map{
			"version":   version,
			"bgColor":   bgColor,
			"fgColor":   fgColor,
			"name":      name,
			"emoji":     emoji,
			"userAgent": ua,
			"startTime": now.Format(dateFormat),
		}
		if userAgentIsCurl(ua) {
			return c.Render("curl", vars)
		}
		return c.Render("index", vars)
	})

	return Webui{
		app: app,
	}
}

// Run starts the Web UI server
func (w Webui) Run() error {
	return w.app.Listen(fmt.Sprintf(":%d", port))
}

func userAgentIsCurl(userAgent string) bool {
	return userAgent[0:4] == "curl"
}
