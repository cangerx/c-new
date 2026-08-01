// Package homepage replaces new-api's default landing page with a custom
// Vue-built homepage, served entirely from this plugin.
//
// Why this works without touching core code: common/embed-file-system.go
// deliberately returns os.ErrNotExist for "/" so the index page falls through
// to gin's NoRoute handler. Because NoRoute only fires when no route matches,
// registering GET "/" here takes precedence — new-api's own index is simply
// never reached. Every other path is untouched, so the React dashboard,
// /api, and /v1 keep working exactly as before.
//
// The frontend lives in frontend/ as an independent Vite project. Rebuild it
// with:
//
//	cd plugins/homepage/frontend && bun install && bun run build
//
// which writes to plugins/homepage/dist/, embedded below.
package homepage

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/pkg/plugin"

	"github.com/gin-gonic/gin"
)

// assetPrefix must not be "/assets" — new-api's static middleware already
// owns that path for its own frontend bundle.
const assetPrefix = "/homepage-assets"

//go:embed all:dist
var distFS embed.FS

// HomepagePlugin serves the custom landing page at "/".
type HomepagePlugin struct {
	assets    fs.FS
	indexPage []byte
}

func (p *HomepagePlugin) Name() string    { return "homepage" }
func (p *HomepagePlugin) Version() string { return "1.0.0" }

// OnInit loads the embedded bundle into memory once, so each request is a
// plain byte-slice write rather than an FS walk.
func (p *HomepagePlugin) OnInit() error {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return err
	}
	p.assets = sub

	index, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		return err
	}
	p.indexPage = index
	common.SysLog("[homepage] custom landing page loaded")
	return nil
}

func (p *HomepagePlugin) RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		// The landing page reflects site settings that an admin can change at
		// any time, so it must not be cached by browsers or proxies.
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", p.indexPage)
	})

	// Vite emits content-hashed filenames, so these are safe to cache hard.
	assetServer := http.StripPrefix(assetPrefix, http.FileServer(http.FS(p.assets)))
	router.GET(assetPrefix+"/*filepath", func(c *gin.Context) {
		if strings.HasSuffix(c.Param("filepath"), "/") {
			c.Status(http.StatusNotFound)
			return
		}
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		assetServer.ServeHTTP(c.Writer, c.Request)
	})

	common.SysLog("[homepage] serving / and " + assetPrefix + "/*")
}

func init() {
	plugin.Register(&HomepagePlugin{})
}
