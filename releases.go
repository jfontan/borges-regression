package regression

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/google/go-github.v15/github"
	"gopkg.in/src-d/go-errors.v1"
)

var releases *Releases

var (
	ErrVersionNotFound = errors.NewKind("Version '%s' not found")
	ErrAssetNotFound   = errors.NewKind(
		"Asset named '%s' not found in release '%s'")
)

type Releases struct {
	client       *github.Client
	repoReleases []*github.RepositoryRelease
}

func GetReleases() *Releases {
	if releases == nil {
		releases = newReleases()
	}

	return releases
}

func newReleases() *Releases {
	return &Releases{
		client: github.NewClient(nil),
	}
}

func (r *Releases) Get(version, asset, path string) error {
	if r.repoReleases == nil {
		err := r.getReleases()
		if err != nil {
			return err
		}
	}

	for _, rel := range r.repoReleases {
		if rel.GetName() == version {
			for _, a := range rel.Assets {
				if a.GetName() == asset {
					return r.download(a.GetBrowserDownloadURL(), path)
				}
			}

			return ErrAssetNotFound.New(asset, version)
		}
	}

	return ErrVersionNotFound.New(version)
}

func (r *Releases) getReleases() error {
	ctx := context.Background()
	rel, _, err := r.client.Repositories.ListReleases(ctx, "src-d", "borges", nil)
	if err != nil {
		return err
	}

	r.repoReleases = rel
	return nil
}

func (r *Releases) download(url, path string) error {
	dir := filepath.Base(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	downloadPath := fmt.Sprintf("%s.download", path)
	exist, err := fileExist(downloadPath)
	if err != nil {
		return err
	}

	if exist {
		err = os.Remove(downloadPath)
		if err != nil {
			return err
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename(downloadPath, path)
	return err
}
