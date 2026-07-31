package plugin

import "github.com/gin-gonic/gin"

// Plugin is the base interface that every plugin must implement.
// It provides metadata used for logging and status reporting.
type Plugin interface {
	// Name returns a unique identifier for this plugin (e.g., "my-auth-hook").
	Name() string

	// Version returns the plugin version (e.g., "1.0.0").
	Version() string
}

// Initializer is an optional interface for plugins that need to run
// initialization logic after all core resources (DB, Redis, config) are ready
// but before the HTTP server starts.
// Return a non-nil error to abort startup.
type Initializer interface {
	Plugin

	// OnInit is called once during InitResources, after model.InitDB(),
	// common.InitRedisClient(), and other core infrastructure is ready.
	OnInit() error
}

// RouteRegistrar is an optional interface for plugins that want to
// register custom HTTP routes and/or middleware.
// Plugins receive the raw *gin.Engine and can add route groups,
// individual routes, or global middleware via engine.Use().
type RouteRegistrar interface {
	Plugin

	// RegisterRoutes is called after all built-in routes have been set up
	// but before the HTTP server starts listening.
	RegisterRoutes(router *gin.Engine)
}

// Shutdowner is an optional interface for plugins that need to perform
// cleanup during graceful shutdown (drain connections, flush buffers, etc.).
type Shutdowner interface {
	Plugin

	// OnShutdown is called during graceful shutdown, before the database
	// connection is closed.
	OnShutdown() error
}
