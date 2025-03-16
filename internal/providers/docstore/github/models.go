package github

// GitHubContent represents the content of a file in a GitHub repository
type GitHubContent struct {
	// Type is the type of content ("file", "dir", "symlink", "submodule")
	Type string `json:"type"`
	// Size is the size of the file in bytes
	Size int `json:"size"`
	// Name is the name of the file
	Name string `json:"name"`
	// Path is the full path to the file in the repository
	Path string `json:"path"`
	// Content is the base64 encoded content of the file
	Content string `json:"content"`
	// Encoding is the encoding of the content (usually "base64")
	Encoding string `json:"encoding"`
	// SHA is the SHA of the file in the repository
	SHA string `json:"sha"`
	// URL is the URL to the file in the GitHub API
	URL string `json:"url"`
	// DownloadURL is the URL to download the file
	DownloadURL string `json:"download_url"`
	// GitURL is the URL to the file in the Git API
	GitURL string `json:"git_url"`
	// HTMLURL is the URL to the file in the GitHub web interface
	HTMLURL string `json:"html_url"`
}

// GitHubFile represents a file to be committed to a GitHub repository
type GitHubFile struct {
	// Path is the path to the file in the repository
	Path string `json:"path"`
	// Content is the content of the file (not base64 encoded)
	Content string `json:"content"`
	// Message is the commit message
	Message string `json:"message"`
	// Branch is the branch to commit to
	Branch string `json:"branch"`
	// SHA is the SHA of the file being replaced (if applicable)
	SHA string `json:"sha,omitempty"`
	// Committer information
	Committer *GitHubCommitter `json:"committer,omitempty"`
}

// GitHubCommitter represents the committer information
type GitHubCommitter struct {
	// Name is the committer's name
	Name string `json:"name"`
	// Email is the committer's email
	Email string `json:"email"`
}

// GitHubCommitResponse represents the response from the GitHub API when committing a file
type GitHubCommitResponse struct {
	// Content is the content of the file
	Content *GitHubContent `json:"content"`
	// Commit is the commit information
	Commit *GitHubCommit `json:"commit"`
}

// GitHubCommit represents a commit in a GitHub repository
type GitHubCommit struct {
	// SHA is the SHA of the commit
	SHA string `json:"sha"`
	// URL is the URL to the commit in the GitHub API
	URL string `json:"url"`
	// HTMLURL is the URL to the commit in the GitHub web interface
	HTMLURL string `json:"html_url"`
}

// GitHubContentListItem represents an item in a list of contents
type GitHubContentListItem struct {
	// Type is the type of content ("file", "dir", "symlink", "submodule")
	Type string `json:"type"`
	// Size is the size of the file in bytes
	Size int `json:"size"`
	// Name is the name of the file
	Name string `json:"name"`
	// Path is the full path to the file in the repository
	Path string `json:"path"`
	// SHA is the SHA of the file in the repository
	SHA string `json:"sha"`
	// URL is the URL to the file in the GitHub API
	URL string `json:"url"`
	// DownloadURL is the URL to download the file
	DownloadURL string `json:"download_url"`
	// GitURL is the URL to the file in the Git API
	GitURL string `json:"git_url"`
	// HTMLURL is the URL to the file in the GitHub web interface
	HTMLURL string `json:"html_url"`
}