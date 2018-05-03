package regression

func NewToolBorges() Tool {
	return Tool{
		Name:        "borges",
		GitURL:      "https://github.com/src-d/borges",
		ProjectPath: "github.com/src-d/borges",
	}
}

func NewBorges(
	config Config,
	version string,
	releases *Releases,
) *Binary {
	return NewBinary(config, NewToolBorges(), version, releases)
}
