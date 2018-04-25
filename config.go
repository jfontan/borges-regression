package regression

import "runtime"

// Config holds the general configuration for tests
type Config struct {
	// Versions has the list of releases to test
	Versions []string
	// OS holds the operating system
	OS string
	// BinaryCache is the path to the borges binaries cache
	BinaryCache string
	// RepositoriesCache is the path to the downloaded repositories
	RepositoriesCache string
	// GitURL is the git repository url to download borges
	GitUrl string
	// GitServerPort is the port where the local git server will listen
	GitServerPort int
}

func NewConfig() Config {
	return Config{
		OS:                runtime.GOOS,
		BinaryCache:       "binary",
		RepositoriesCache: "repos",
		GitUrl:            "https://github.com/src-d/borges",
		GitServerPort:     9418,
	}
}
