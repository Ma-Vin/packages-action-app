package github_model

type JsonPackageType string
type JsonVisibility string
type JsonActivation string

const (
	NPM       JsonPackageType = "npm"
	MAVEN     JsonPackageType = "maven"
	RUBYGEMS  JsonPackageType = "rubygems"
	DOCKER    JsonPackageType = "docker"
	NUGET     JsonPackageType = "nuget"
	CONTAINER JsonPackageType = "container"
)

const (
	PRIVATE JsonVisibility = "private"
	PUBLIC  JsonVisibility = "public"
)

const (
	ENABLED  JsonActivation = "enabled"
	DISABLED JsonActivation = "disabled"
)
