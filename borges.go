package regression

func NewBorges(config Config, version string, releases *Releases) *Binary {
	return NewBinary(config, "borges", version, releases)
}
