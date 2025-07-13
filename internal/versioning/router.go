// Package versioning provides utilities for API version management.
package versioning

import (
	"github.com/gin-gonic/gin"
)

// Router is a helper for managing versioned API routes
type Router struct {
	engine *gin.Engine
	groups map[Version]*gin.RouterGroup
}

// NewRouter creates a new versioned router
func NewRouter(engine *gin.Engine) *Router {
	return &Router{
		engine: engine,
		groups: make(map[Version]*gin.RouterGroup),
	}
}

// Group returns a router group for the specified version
// If the version doesn't exist, it creates a new group
func (r *Router) Group(version Version) *gin.RouterGroup {
	if group, exists := r.groups[version]; exists {
		return group
	}
	
	// Create a new group for this version
	group := r.engine.Group(GetPathPrefix(version))
	r.groups[version] = group
	return group
}

// Latest returns a router group for the latest API version
func (r *Router) Latest() *gin.RouterGroup {
	return r.Group(LatestVersion)
}

// RegisterVersionedEndpoint registers an endpoint across multiple API versions
// This is useful when an endpoint remains the same across versions
func (r *Router) RegisterVersionedEndpoint(path string, method string, versions []Version, handler gin.HandlerFunc) {
	for _, version := range versions {
		group := r.Group(version)
		
		switch method {
		case "GET":
			group.GET(path, handler)
		case "POST":
			group.POST(path, handler)
		case "PUT":
			group.PUT(path, handler)
		case "DELETE":
			group.DELETE(path, handler)
		case "PATCH":
			group.PATCH(path, handler)
		case "HEAD":
			group.HEAD(path, handler)
		case "OPTIONS":
			group.OPTIONS(path, handler)
		}
	}
}

// RegisterAllVersions registers an endpoint across all supported API versions
func (r *Router) RegisterAllVersions(path string, method string, handler gin.HandlerFunc) {
	r.RegisterVersionedEndpoint(path, method, All(), handler)
}