package regression

import (
	"fmt"
	"os"
	"path/filepath"

	log "gopkg.in/src-d/go-log.v0"
)

type RepoDescription struct {
	Name        string
	URL         string
	Description string
	Complexity  int
}

var defaultRepos = []RepoDescription{
	{
		Name:        "cangallo",
		URL:         "git://github.com/jfontan/cangallo.git",
		Description: "Small repository that should be fast to clone",
		Complexity:  0,
	},
	{
		Name:        "octoprint-tft",
		URL:         "https://github.com/mcuadros/OctoPrint-TFT",
		Description: "Small repository that should be fast to clone",
		Complexity:  0,
	},
	{
		Name:        "upsilon",
		URL:         "git://github.com/upsilonproject/upsilon-common.git",
		Description: "Average repository",
		Complexity:  1,
	},
	{
		Name:        "numpy",
		URL:         "git://github.com/numpy/numpy.git",
		Description: "Average repository",
		Complexity:  2,
	},
	{
		Name:        "tensorflow",
		URL:         "git://github.com/tensorflow/tensorflow.git",
		Description: "Average repository",
		Complexity:  3,
	},
	{
		Name:        "bismuth",
		URL:         "git://github.com/hclivess/Bismuth.git",
		Description: "Big files repo (100Mb)",
		Complexity:  4,
	},
}

type Repositories struct {
	repos  []RepoDescription
	config Config
}

func NewRepositories(config Config) *Repositories {
	return &Repositories{
		repos:  defaultRepos,
		config: config,
	}
}

func (r *Repositories) Download() error {
	for _, repo := range r.repos {
		if repo.Complexity > r.config.Complexity {
			continue
		}

		logger, _ := log.New()
		logger = logger.New(log.Fields{"name": repo.Name})

		path := filepath.Join(r.config.RepositoriesCache, repo.Name)
		exist, err := fileExist(path)
		if err != nil {
			return err
		}
		if exist {
			logger.Debugf("Repository already downloaded")
			continue
		}

		logger = logger.New(log.Fields{
			"url":  repo.URL,
			"path": path,
		})

		logger.Debugf("Downloading repository")
		err = os.MkdirAll(r.config.RepositoriesCache, 0755)
		if err != nil {
			return err
		}

		err = downloadRepo(logger, repo.URL, path)
		if err != nil {
			logger.Error(err, "Could not download repository")
			return err
		}
	}

	return nil
}

func (r *Repositories) Path() string {
	return r.config.RepositoriesCache
}

func (r *Repositories) Names(complexity int) []string {
	names := make([]string, 0, len(r.repos))
	for _, repo := range r.repos {
		if repo.Complexity <= complexity {
			names = append(names, repo.Name)
		}
	}

	return names
}

func downloadRepo(l log.Logger, url, path string) error {
	downloadPath := fmt.Sprintf("%s.download", path)
	exist, err := fileExist(downloadPath)
	if err != nil {
		return err
	}

	if exist {
		err = os.RemoveAll(downloadPath)
		if err != nil {
			return err
		}
	}

	clone, err := NewExecutor("git", "clone", "--bare", url, downloadPath)
	if err != nil {
		l.Error(err, "Could not create executor")
		return err
	}

	err = clone.Run()
	if err != nil {
		out, _ := clone.Out()
		l.New(log.Fields{"output": out}).Error(err, "Could not execute git clone")
		return err
	}

	err = os.Rename(downloadPath, path)
	return err
}
