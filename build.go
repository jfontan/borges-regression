package regression

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/src-d/go-errors.v0"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Build struct {
	// Version is the reference that will be built
	Version string

	// GoPath is the directory where the temporary path where borges is built
	GoPath string

	source    string
	reference string
	url       string
}

var borgesRepo = "https://github.com/src-d/borges"
var regRepo = regexp.MustCompile(`^(local|remote|pull):([[:ascii:]]+)$`)

var (
	ErrReferenceNotFound = errors.NewKind("Reference %s not found")
	ErrInvalidVersion    = errors.NewKind("Version %s is invalid")
)

func IsRepo(version string) bool {
	return regRepo.MatchString(version)
}

func NewBuild(version string) (*Build, error) {
	if !IsRepo(version) {
		return nil, ErrInvalidVersion.New(version)
	}

	source, reference := parseVersion(version)

	url := borgesRepo
	if source == "local" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		url = fmt.Sprintf("file://%s", pwd)
	}

	return &Build{
		Version:   version,
		source:    source,
		reference: reference,
		url:       url,
	}, nil
}

func (b *Build) download() error {
	dir, err := createTempDir()
	if err != nil {
		return err
	}

	b.GoPath = dir

	clonePath := filepath.Join(dir, "src", "github.com", "src-d", "borges")
	err = os.MkdirAll(clonePath, 0755)
	if err != nil {
		return err
	}

	r, err := git.PlainInit(clonePath, false)
	if err != nil {
		return err
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{b.url},
	})

	referenceName, err := findReference(b.Version, remote)
	if err != nil {
		return err
	}

	refSpecs := []config.RefSpec{
		config.RefSpec(fmt.Sprintf("%s:refs/heads/master", referenceName)),
	}

	err = r.Fetch(&git.FetchOptions{
		Depth:    1,
		RefSpecs: refSpecs,
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/master"),
	})

	return err
}

func findReference(
	version string,
	remote *git.Remote,
) (string, error) {
	source, reference := parseVersion(version)
	if source == "pull" {
		return fmt.Sprintf("refs/pull/%s/head", reference), nil
	}

	refs, err := remote.List(new(git.ListOptions))
	if err != nil {
		return "", err
	}

	for _, ref := range refs {
		name := ref.Name()

		if name.IsBranch() || name.IsTag() {
			if name.Short() == reference {
				return name.String(), nil
			}
		}
	}

	return "", ErrReferenceNotFound.New(reference)
}

func parseVersion(version string) (string, string) {
	r := regRepo.FindStringSubmatch(version)
	spew.Dump(r)
	return r[1], r[2]
}
