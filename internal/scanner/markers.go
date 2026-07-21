package scanner

import "strings"

// wellKnownBasenames are marker files detected by basename anywhere in the tree.
var wellKnownBasenames = map[string]string{
	"package.json":      "package.json",
	"package-lock.json": "package-lock.json",
	"pnpm-lock.yaml":    "pnpm-lock.yaml",
	"yarn.lock":         "yarn.lock",
	"tsconfig.json":     "tsconfig.json",
	"nest-cli.json":     "nest-cli.json",
	"requirements.txt":  "requirements.txt",
	"pyproject.toml":    "pyproject.toml",
	"poetry.lock":       "poetry.lock",
	"uv.lock":           "uv.lock",
	"go.mod":            "go.mod",
	"cargo.toml":        "Cargo.toml",
	"pom.xml":           "pom.xml",
	"dockerfile":        "Dockerfile",
	"docker-compose.yml": "docker-compose.yml",
	"docker-compose.yaml": "docker-compose.yaml",
	"compose.yaml":      "compose.yaml",
	"compose.yml":       "compose.yml",
	"schema.prisma":     "schema.prisma",
}

func recordMarkers(dst map[string]struct{}, relSlash, name string) {
	lowerName := strings.ToLower(name)
	if marker, ok := wellKnownBasenames[lowerName]; ok {
		dst[marker] = struct{}{}
	}

	// Prefix / glob style markers.
	switch {
	case strings.HasPrefix(lowerName, "next.config."):
		dst["next.config.*"] = struct{}{}
	case strings.HasPrefix(lowerName, "vite.config."):
		dst["vite.config.*"] = struct{}{}
	case strings.HasPrefix(lowerName, "webpack.config."):
		dst["webpack.config.*"] = struct{}{}
	case strings.HasPrefix(lowerName, "build.gradle"):
		dst["build.gradle"] = struct{}{}
	case strings.HasPrefix(lowerName, "dockerfile"):
		dst["Dockerfile"] = struct{}{}
	}

	relLower := strings.ToLower(relSlash)
	switch {
	case strings.HasPrefix(relLower, ".github/workflows/"):
		dst[".github/workflows"] = struct{}{}
	case strings.HasPrefix(relLower, "helm/") || strings.Contains(relLower, "/helm/"):
		dst["helm/"] = struct{}{}
	case strings.HasPrefix(relLower, "sql/") || strings.Contains(relLower, "/sql/"):
		dst["sql/"] = struct{}{}
	case strings.Contains(relLower, "/migrations/") || strings.HasPrefix(relLower, "migrations/"):
		dst["migrations/"] = struct{}{}
	case strings.HasSuffix(relLower, "/schema.prisma") || relLower == "schema.prisma" || strings.HasSuffix(relLower, "prisma/schema.prisma"):
		dst["schema.prisma"] = struct{}{}
	case relLower == ".gitignore" || strings.HasSuffix(relLower, "/.gitignore"):
		dst[".gitignore"] = struct{}{}
	}
}
