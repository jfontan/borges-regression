package regression

import "runtime"

// Config holds the general configuration for tests
type Config struct {
	// Versions has the list of releases to test
	Versions []string
	// OS holds the operating system
	OS string
	// BinaryCache is the path to the borges binaries cache
	BinaryCache string `env:"REG_BINARIES" default:"binaries" long:"binaries" description:"Directory to store borges binaries"`
	// RepositoriesCache is the path to the downloaded repositories
	RepositoriesCache string `env:"REG_REPOS" default:"repos" long:"repos" description:"Directory to store repositories"`
	// GitURL is the git repository url to download borges
	GitURL string `env:"REG_GITURL" default:"https://github.com/src-d/borges" long:"url" description:"URL to borges repo"`
	// GitServerPort is the port where the local git server will listen
	GitServerPort int `env:"REG_GITPORT" default:"9418" long:"gitport" description:"Port for local git server"`
	// Complexity has the max number of complexity of repos to test
	Complexity int `env:"REG_COMPLEXITY" default:"1" long:"complexity" short:"c" description:"Complexity of the repositories to test"`
	// Repeat is the number of times each test will be run
	Repeat int `env:"REG_REPEAT" default:"3" long:"repeat" short:"n" description:"Number of times a test is run"`
	// ShowRepos when --show-repos is specified
	ShowRepos bool `long:"show-repos" description:"List available repositories to test"`
}

func NewConfig() Config {
	return Config{
		OS: runtime.GOOS,
	}
}
