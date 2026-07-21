package domain

// Category classifies a finding by concern area.
type Category string

const (
	CategoryLanguage  Category = "language"
	CategoryFramework Category = "framework"
	CategoryInfra     Category = "infra"
	CategoryDatabase  Category = "database"
	CategoryArch      Category = "architecture"
	CategoryDocs      Category = "documentation"
	CategoryVCS       Category = "vcs"
	CategoryCI        Category = "ci"
	CategoryTooling   Category = "tooling"
)
