package plugins

import (
	"github.com/FiaLDI/project-parse/internal/plugins/architecture"
	"github.com/FiaLDI/project-parse/internal/plugins/database"
	"github.com/FiaLDI/project-parse/internal/plugins/docker"
	docplugin "github.com/FiaLDI/project-parse/internal/plugins/documentation"
	"github.com/FiaLDI/project-parse/internal/plugins/git"
	"github.com/FiaLDI/project-parse/internal/plugins/githubactions"
	"github.com/FiaLDI/project-parse/internal/plugins/golang"
	"github.com/FiaLDI/project-parse/internal/plugins/java"
	"github.com/FiaLDI/project-parse/internal/plugins/kubernetes"
	"github.com/FiaLDI/project-parse/internal/plugins/node"
	"github.com/FiaLDI/project-parse/internal/plugins/python"
	"github.com/FiaLDI/project-parse/internal/plugins/rust"
	"github.com/FiaLDI/project-parse/internal/ports"
)

// RegisterAll registers built-in analysis plugins.
func RegisterAll(r ports.Registry, cache ports.FileCache) {
	if r == nil {
		return
	}
	for _, p := range []ports.Plugin{
		node.New(cache),
		python.New(cache),
		golang.New(cache),
		rust.New(cache),
		java.New(cache),
		docker.New(cache),
		git.New(cache),
		githubactions.New(cache),
		kubernetes.New(cache),
		database.New(cache),
		architecture.New(cache),
		docplugin.New(cache),
	} {
		r.Register(p)
	}
}
