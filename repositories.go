package regression

import (
	"fmt"
	"os"
	"path/filepath"
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
		Name:        "upsilon",
		URL:         "git://github.com/upsilonproject/upsilon-common.git",
		Description: "Average repository",
		Complexity:  1,
	},
	// {
	// 	Name:        "numpy",
	// 	URL:         "git://github.com/numpy/numpy.git",
	// 	Description: "Average repository",
	// 	Complexity:  2,
	// },
	// {
	// 	Name:        "tensorflow",
	// 	URL:         "git://github.com/tensorflow/tensorflow.git",
	// 	Description: "Average repository",
	// 	Complexity:  3,
	// },
	// {
	// 	Name:        "bismuth",
	// 	URL:         "git://github.com/hclivess/Bismuth.git",
	// 	Description: "Big files repo (100Mb)",
	// 	Complexity:  4,
	// },
}

type Repositories struct {
	repos    []RepoDescription
	cacheDir string
}

func NewRepositories() *Repositories {
	return &Repositories{
		repos:    defaultRepos,
		cacheDir: "repos",
	}
}

func (r *Repositories) Download() error {
	for _, repo := range r.repos {
		path := filepath.Join(r.cacheDir, repo.Name)
		exist, err := fileExist(path)
		if err != nil {
			return err
		}
		if exist {
			println("Repository", repo.Name, "already in cache")
			continue
		}

		println("Downloading", repo.Name)
		err = os.MkdirAll(r.cacheDir, 0755)
		if err != nil {
			return err
		}

		err = downloadRepo(repo.URL, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repositories) Path() string {
	return r.cacheDir
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

func downloadRepo(url, path string) error {
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
		// TODO: log error
		return err
	}

	err = clone.Run()
	if err != nil {
		// TODO: log error output
		return err
	}

	err = os.Rename(downloadPath, path)
	return err
}
