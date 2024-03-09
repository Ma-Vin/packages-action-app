package github_model

// package definition for a user, see also: https://docs.github.com/en/rest/packages/packages?apiVersion=2022-11-28#get-a-package-for-a-user
type UserPackage struct {
	Id           int             `json:"id"`
	Name         string          `json:"name"`
	PackageType  JsonPackageType `json:"package_type"`
	Url          string          `json:"url"`
	HtmlUrl      string          `json:"html_url"`
	VersionCount int             `json:"version_count"`
	Visibility   JsonVisibility  `json:"visibility"`
	Owner        User            `json:"owner"`
	Repository   Repository      `json:"repository"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type User struct {
	Name              string `json:"name"`
	EMail             string `json:"email"`
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistUrl           string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredAt         string `json:"starred_at"`
}

type Repository struct {
	Id                       int                 `json:"id"`
	NodeID                   string              `json:"node_id"`
	Name                     string              `json:"name"`
	FullName                 string              `json:"full_name"`
	Owner                    User                `json:"owner"`
	Private                  bool                `json:"private"`
	HtmlUrl                  string              `json:"html_url"`
	Description              string              `json:"description"`
	Fork                     bool                `json:"fork"`
	Url                      string              `json:"url"`
	ArchiveUrl               string              `json:"archive_url"`
	AssigneesUrl             string              `json:"assignees_url"`
	BlobsUrl                 string              `json:"blobs_url"`
	BranchesUrl              string              `json:"branches_url"`
	CollaboratorsUrl         string              `json:"collaborators_url"`
	CommentsUrl              string              `json:"comments_url"`
	CommitsUrl               string              `json:"commits_url"`
	CompareUrl               string              `json:"compare_url"`
	ContentsUrl              string              `json:"contents_url"`
	ContributorsUrl          string              `json:"contributors_url"`
	DeploymentsUrl           string              `json:"deployments_url"`
	DownloadsUrl             string              `json:"downloads_url"`
	EventsUrl                string              `json:"events_url"`
	ForksUrl                 string              `json:"forks_url"`
	GitCommitsUrl            string              `json:"git_commits_url"`
	GitRefsUrl               string              `json:"git_refs_url"`
	GitTagsUrl               string              `json:"git_tags_url"`
	GitUrl                   string              `json:"git_url"`
	IssueCommentUrl          string              `json:"issue_comment_url"`
	IssueEventsUrl           string              `json:"issue_events_url"`
	IssuesUrl                string              `json:"issues_url"`
	KeysUrl                  string              `json:"keys_url"`
	LabelsUrl                string              `json:"labels_url"`
	LanguagesUrl             string              `json:"languages_url"`
	MergesUrl                string              `json:"merges_url"`
	MilestonesUrl            string              `json:"milestones_url"`
	NotificationsUrl         string              `json:"notifications_url"`
	PullsUrl                 string              `json:"pulls_url"`
	ReleasesUrl              string              `json:"releases_url"`
	SshUrl                   string              `json:"ssh_url"`
	StargazersUrl            string              `json:"stargazers_url"`
	StatusesUrl              string              `json:"statuses_url"`
	SubscribersUrl           string              `json:"subscribers_url"`
	SubscriptionUrl          string              `json:"subscription_url"`
	TagsUrl                  string              `json:"tags_url"`
	TeamsUrl                 string              `json:"teams_url"`
	TreesUrl                 string              `json:"trees_url"`
	CloneUrl                 string              `json:"clone_url"`
	MirrorUrl                string              `json:"mirror_url"`
	HooksUrl                 string              `json:"hooks_url"`
	SvnUrl                   string              `json:"svn_url"`
	Homepage                 string              `json:"homepage"`
	Language                 string              `json:"language"`
	ForksCount               int                 `json:"forks_count"`
	StargazersCount          int                 `json:"stargazers_count"`
	WatchersCount            int                 `json:"watchers_count"`
	Size                     int                 `json:"size"`
	DefaultBranch            string              `json:"default_branch"`
	OpenIssuesCount          int                 `json:"open_issues_count"`
	IsTemplate               bool                `json:"is_template"`
	Topics                   []string            `json:"topics"`
	HasIssues                bool                `json:"has_issues"`
	HasProjects              bool                `json:"has_projects"`
	HasWiki                  bool                `json:"has_wiki"`
	HasPages                 bool                `json:"has_pages"`
	HasDownloads             bool                `json:"has_downloads"`
	HasDiscussions           bool                `json:"has_discussions"`
	Archived                 bool                `json:"archived"`
	Visibility               JsonVisibility      `json:"visibility"`
	PushedAt                 string              `json:"pushed_at"`
	CreatedAt                string              `json:"created_at"`
	UpdatedAt                string              `json:"updated_at"`
	Permissions              Permissons          `json:"permissions"`
	RoleName                 string              `json:"role_name"`
	TempCloneToken           string              `json:"temp_clone_token"`
	SubscribersCount         int                 `json:"subscribers_count"`
	NetworkCount             int                 `json:"network_count"`
	CodeOfConduct            CodeOfConduct       `json:"code_of_conduct"`
	License                  License             `json:"license"`
	Forks                    int                 `json:"forks"`
	OpenIssues               int                 `json:"open_issues"`
	Watchers                 int                 `json:"watchers"`
	AllowForking             bool                `json:"allow_forking"`
	WebCommitSignoffRequired bool                `json:"web_commit_signoff_required"`
	SecurityAndAnalysis      SecurityAndAnalysis `json:"security_and_analysis"`
}

type Permissons struct {
	Admin    bool `json:"admin"`
	Maintain bool `json:"maintain"`
	Push     bool `json:"push"`
	Triage   bool `json:"triage"`
	Pull     bool `json:"pull"`
}

type CodeOfConduct struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Body    string `json:"body"`
	HtmlUrl string `json:"html_url"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxId string `json:"spdx_id"`
	Url    string `json:"url"`
	NodeId string `json:"node_id"`
}

type SecurityAndAnalysis struct {
	AdvancedSecurity             JsonActivation `json:"advanced_security"`
	DependabotSecurityUpdates    JsonActivation `json:"dependabot_security_updates"`
	SecretScanning               JsonActivation `json:"secret_scanning"`
	SecretScanningPushProtection JsonActivation `json:"secret_scanning_push_protection"`
}
