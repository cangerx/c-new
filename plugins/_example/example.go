// Package example demonstrates how to write a plugin for new-api.
// It implements all four plugin interfaces: Plugin, Initializer,
// RouteRegistrar, and Shutdowner.
//
// Use this as a template for your own plugins.
package example

import (
	"fmt"
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/plugin"

	"github.com/gin-gonic/gin"
)

// ExamplePlugin demonstrates plugin capabilities.
type ExamplePlugin struct{}

func (p *ExamplePlugin) Name() string    { return "example" }
func (p *ExamplePlugin) Version() string { return "1.0.0" }

// OnInit is called after all core infrastructure (DB, Redis, etc.) is ready.
// Use it for config validation, DB migrations, starting background tasks, etc.
func (p *ExamplePlugin) OnInit() error {
	common.SysLog("[example-plugin] OnInit called")

	// Demonstrate DB access via model package
	var userCount int64
	if err := model.DB.Raw("SELECT COUNT(*) FROM users").Scan(&userCount).Error; err != nil {
		return fmt.Errorf("failed to query DB: %w", err)
	}
	common.SysLog(fmt.Sprintf("[example-plugin] user count: %d", userCount))
	return nil
}

// RegisterRoutes is called after all built-in routes are set up.
// Plugins can add new route groups, individual routes, or global middleware.
func (p *ExamplePlugin) RegisterRoutes(router *gin.Engine) {
	// Public status endpoint (no auth)
	router.GET("/plugin/example/status", func(c *gin.Context) {
		plugins := plugin.GetAll()
		info := make([]gin.H, len(plugins))
		for i, pl := range plugins {
			info[i] = gin.H{
				"name":    pl.Name(),
				"version": pl.Version(),
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"plugins": info,
		})
	})
}

// OnShutdown is called during graceful shutdown.
// Use it to drain connections, flush buffers, or clean up resources.
func (p *ExamplePlugin) OnShutdown() error {
	common.SysLog("[example-plugin] OnShutdown called")
	return nil
}

func init() {
	plugin.Register(&ExamplePlugin{})
}
