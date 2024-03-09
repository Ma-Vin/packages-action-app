package github_model

// package version definition for a user, see also: https://docs.github.com/en/rest/packages/packages?apiVersion=2022-11-28#get-a-package-version-for-a-user
type Version struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Url            string   `json:"url"`
	PackageHtmlUrl string   `json:"package_html_url"`
	HtmlUrl        string   `json:"html_url"`
	License        string   `json:"license"`
	Description    string   `json:"description"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
	DeletedAt      string   `json:"deleted_at"`
	Metadata       Metadata `json:"metadata"`
}

type Metadata struct {
	PackageType JsonPackageType `json:"package_type"`
	Container   Container       `json:"container"`
	Docker      Docker          `json:"docker"`
}

type Container struct {
	Tags []string `json:"tags"`
}

type Docker struct {
	Tags []string `json:"tag"`
}
