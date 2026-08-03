// Package plugins aggregates all active plugins through blank imports.
// Plugins register themselves via their init() functions.
//
// To enable a plugin, add its import path here. To disable, comment it out.
//
//	import _ "github.com/QuantumNous/new-api/plugins/my-plugin"
//
// The main package imports this package to auto-register all active plugins.
package plugins

import (
	// Example plugin — demonstrates plugin capabilities.
	// Remove or comment out this line to disable the example.
	_ "github.com/QuantumNous/new-api/plugins/_example"

	// Custom landing page, replaces the default index at "/".
	// Disabled: the landing page now lives in the React app at
	// web/src/features/landing, wired to "/" via web/src/routes/index.tsx.
	// Re-enabling this would shadow that route, because an explicit Gin
	// GET("/") outranks the NoRoute handler that serves the SPA.
	// _ "github.com/QuantumNous/new-api/plugins/homepage"
)
