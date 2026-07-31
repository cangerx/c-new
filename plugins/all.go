// Package plugins aggregates all active plugins through blank imports.
// Plugins register themselves via their init() functions.
//
// To enable a plugin, add its import path here. To disable, comment it out.
//
//	import _ "github.com/QuantumNous/new-api/plugins/my-plugin"
//
// The main package imports this package to auto-register all active plugins.
package plugins

// Example plugin — demonstrates plugin capabilities.
// Remove or comment out this line to disable the example.
import _ "github.com/QuantumNous/new-api/plugins/_example"
