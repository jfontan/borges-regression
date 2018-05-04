package regression

func NewToolGitbase() Tool {
	return Tool{
		Name:        "gitbase",
		GitURL:      "https://github.com/src-d/gitbase",
		ProjectPath: "github.com/src-d/gitbase",
	}
}

func NewGitbase(
	config Config,
	version string,
	releases *Releases,
) *Binary {
	return NewBinary(config, NewToolGitbase(), version, releases)
}
