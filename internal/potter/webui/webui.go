package webui

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JulienBreux/potter/pkg/color"
	"github.com/JulienBreux/potter/pkg/emoji"
	"github.com/JulienBreux/potter/pkg/namesgen"
	"github.com/JulienBreux/potter/pkg/version"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/inhies/go-bytesize"
	"github.com/pbnjay/memory"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// TODO: Move to configuration
	port = 8080

	dateFormat = "02-01-2006 15:04:05"

	directory = "./views"
	extension = ".go.html"
)

// WebUI represents a Web UI server
type Webui struct {
	app *fiber.App
}

// New creates a new Web UI server
func New() Webui {
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
			"version":   version.Version,
			"commit":    version.Commit,
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
	app.Get("/memory", pageMemory)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		// TODO: Log
		return c.Status(http.StatusOK).SendString("ok")
	})
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

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

// TODO: Move to ignore file + function
var ignoreFilesPrefixes = [...]string{
	// "/app",
	"/.git", // "/.dockerenv",
	"/bin", "/boot", "/dev", "/dev", "/etc",
	"/lib", "/proc", "/root", "/run", "/sys",
	"/usr", "/var", "/views", "/sbin",
}

func pageFiles(c *fiber.Ctx) error {
	files := make(map[string]bytesize.ByteSize)

	err := filepath.WalkDir("/",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			for _, prefix := range ignoreFilesPrefixes {
				if strings.HasPrefix(path, prefix) {
					return filepath.SkipDir
				}
			}

			if path == "/" && d.IsDir() {
				return nil
			}

			var bs bytesize.ByteSize
			if f, err := d.Info(); err == nil {
				bs = bytesize.New(float64(f.Size()))
			}
			files[path] = bs
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
	path := c.Query("path")

	// TODO: Move
	for _, prefix := range ignoreFilesPrefixes {
		if strings.HasPrefix(path, prefix) {
			return c.Redirect("/files?updateError=1&path=" + path + "&error=Path forbidden")
		}
	}

	newPath := false
	content, err := os.ReadFile(path)
	if err != nil {
		newPath = true
	}

	vars := fiber.Map{
		"path":    path,
		"content": string(content),
		"new":     newPath,
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

	// TODO: Move
	for _, prefix := range ignoreFilesPrefixes {
		if strings.HasPrefix(path, prefix) {
			return c.Redirect("/files?deleteError=1&path=" + path + "&error=Path forbidden")
		}
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
		return c.Redirect("/files?updateError=1&path=" + f.Path + "&error=" + err.Error())
	}

	return c.Redirect("/files?updateSuccess=1&path=" + f.Path)
}

func pageMemory(c *fiber.Ctx) error {
	vars := fiber.Map{
		"memoryCurrent": "todo",
		"memoryFree":    bytesize.New(float64(memory.FreeMemory())),
		"memoryTotal":   bytesize.New(float64(memory.TotalMemory())),
	}
	return c.Render("memory", vars)
}

func normalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13, 10}, []byte{10})
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.ReplaceAll(d, []byte{13}, []byte{10})
	return d
}
