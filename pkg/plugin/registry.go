package plugin

import (
	"fmt"
	"sync"

	"github.com/QuantumNous/new-api/common"

	"github.com/gin-gonic/gin"
)

var (
	registry   []Plugin
	registryMu sync.Mutex
)

// Register adds a plugin to the global registry. Plugins call this from their
// init() function. Registration order determines initialization order.
// Panics if a plugin with the same name is already registered.
func Register(p Plugin) {
	registryMu.Lock()
	defer registryMu.Unlock()

	name := p.Name()
	for _, existing := range registry {
		if existing.Name() == name {
			panic(fmt.Sprintf("plugin %q is already registered", name))
		}
	}
	registry = append(registry, p)
}

// GetAll returns all registered plugins in registration order.
func GetAll() []Plugin {
	registryMu.Lock()
	defer registryMu.Unlock()

	result := make([]Plugin, len(registry))
	copy(result, registry)
	return result
}

// InitAll calls OnInit on every registered plugin that implements Initializer.
// Stops on the first error and returns it.
// Must be called after model.InitDB() and other core infrastructure is ready.
func InitAll() error {
	plugins := GetAll()
	for _, p := range plugins {
		init, ok := p.(Initializer)
		if !ok {
			continue
		}
		common.SysLog(fmt.Sprintf("plugin %q initializing...", p.Name()))
		if err := init.OnInit(); err != nil {
			common.SysError(fmt.Sprintf("plugin %q init failed: %v", p.Name(), err))
			return fmt.Errorf("plugin %q: %w", p.Name(), err)
		}
		common.SysLog(fmt.Sprintf("plugin %q initialized", p.Name()))
	}
	return nil
}

// RegisterAllRoutes calls RegisterRoutes on every registered plugin that
// implements RouteRegistrar, passing the gin.Engine.
// Must be called after router.SetRouter() but before the server starts listening.
func RegisterAllRoutes(router *gin.Engine) {
	plugins := GetAll()
	for _, p := range plugins {
		rp, ok := p.(RouteRegistrar)
		if !ok {
			continue
		}
		rp.RegisterRoutes(router)
		common.SysLog(fmt.Sprintf("plugin %q routes registered", p.Name()))
	}
}

// ShutdownAll calls OnShutdown on every registered plugin that implements
// Shutdowner, in reverse registration order (LIFO).
// Errors are collected; all plugins get a chance to shut down.
func ShutdownAll() error {
	plugins := GetAll()
	var errs []error
	for i := len(plugins) - 1; i >= 0; i-- {
		p := plugins[i]
		sd, ok := p.(Shutdowner)
		if !ok {
			continue
		}
		common.SysLog(fmt.Sprintf("plugin %q shutting down...", p.Name()))
		if err := sd.OnShutdown(); err != nil {
			common.SysError(fmt.Sprintf("plugin %q shutdown error: %v", p.Name(), err))
			errs = append(errs, fmt.Errorf("plugin %q: %w", p.Name(), err))
		} else {
			common.SysLog(fmt.Sprintf("plugin %q shut down", p.Name()))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}
