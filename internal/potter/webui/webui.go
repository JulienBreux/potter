package webui

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JulienBreux/potter/pkg/color"
	"github.com/JulienBreux/potter/pkg/emoji"
	"github.com/JulienBreux/potter/pkg/namesgen"
	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
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

		Views:       engine,
		ViewsLayout: "layouts/main",
	}
	app := fiber.New(config)

	// TODO: Move to config
	color := color.Rand()
	bgColor := color.ToHex()
	fgColor := color.Invert().ToHex()
	name := namesgen.GetRandom()
	emoji := emoji.GetRandom()

	app.Static("assets", "./assets")

	app.Use(func(c *fiber.Ctx) error {
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
		if err := c.Bind(vars); err != nil {
			return err
		}
		return c.Next()
	})

	app.Get("/", pageIndex)
	app.Get("/vars", pageVars)
	app.Get("/files", pageFiles)
	app.Get("/files/read", pageFilesRead)
	app.Post("/files/update", pageFilesUpdate)
	app.Get("/files/delete", pageFilesDelete)

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

func pageIndex(c *fiber.Ctx) error {
	if userAgentIsCurl(string(c.Context().UserAgent())) {
		return c.Render("curl/index", fiber.Map{}, "")
	}
	return c.Render("index", fiber.Map{})
}

func pageVars(c *fiber.Ctx) error {
	envvars := make(map[string]string)
	// TODO: Move to package
	for _, keyval := range os.Environ() {
		v := strings.Split(keyval, "=")
		envvars[v[0]] = v[1]
	}
	vars := fiber.Map{"vars": envvars}
	if userAgentIsCurl(string(c.Context().UserAgent())) {
		return c.Render("curl/vars", vars, "")
	}
	return c.Render("vars", vars)
}

func pageFiles(c *fiber.Ctx) error {
	files := make(map[string]string)

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == "." {
				return nil
			}
			files[path] = humanize.Bytes(uint64(info.Size()))
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	vars := fiber.Map{
		"files": files,

		"path":          c.Query("path"),
		"error":         c.Query("error"),
		"deleteSuccess": c.Query("deleteSuccess"),
		"deleteError":   c.Query("deleteError"),
		"updateSuccess": c.Query("updateSuccess"),
		"updateError":   c.Query("updateError"),
	}
	if userAgentIsCurl(string(c.Context().UserAgent())) {
		return c.Render("curl/files", vars, "")
	}
	return c.Render("files/files", vars)
}

func pageFilesRead(c *fiber.Ctx) error {
	file := c.Query("path")

	newFile := false
	content, err := os.ReadFile(file)
	if err != nil {
		newFile = true
	}

	vars := fiber.Map{
		"path":    c.Query("path"),
		"content": string(content),
		"new":     newFile,
	}
	return c.Render("files/read", vars)
}

func pageFilesDelete(c *fiber.Ctx) error {
	path := c.Query("path")

	// TODO: Move action
	if strings.HasPrefix(path, "/") ||
		strings.HasPrefix(path, "./") ||
		strings.HasPrefix(path, "../") {
		return c.Redirect("/files?deleteError=1&path=" + path)
	}

	if err := os.Remove(path); err != nil {
		return c.Redirect("/files?deleteError=1&path=" + path + "&error=" + err.Error())
	}

	return c.Redirect("/files?deleteSuccess=1&path=" + path)
}

type file struct {
	Path, Content string
}

func pageFilesUpdate(c *fiber.Ctx) error {
	f := new(file)

	if err := c.BodyParser(f); err != nil {
		return err
	}

	// TODO: Define perms
	ctnt := normalizeNewlines([]byte(f.Content))
	const perm = 0600
	if err := os.WriteFile(f.Path, ctnt, perm); err != nil {
		return err
	}

	return c.Redirect("/files?updateSuccess=1&path=" + f.Path)
}

func normalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13, 10}, []byte{10})
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13}, []byte{10})
	return d
}
