package regression

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var regRelease = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

type Borges struct {
	Version string
	Distro  string
	Path    string

	binCache string
}

func NewBorges(version string) *Borges {
	return &Borges{
		Version:  version,
		Distro:   "linux",
		binCache: "binaries",
	}
}

func (b *Borges) IsRelease() bool {
	return regRelease.MatchString(b.Version)
}

func (b *Borges) Download() error {
	if !b.IsRelease() {
		b.Path = b.Version
		return nil
	}

	cacheName := b.cacheName()
	exist, err := fileExist(cacheName)
	if err != nil {
		return err
	}

	if exist {
		b.Path = cacheName
		return nil
	}

	err = b.downloadRelease()
	if err != nil {
		return err
	}

	b.Path = cacheName

	return nil
}

func (b *Borges) downloadRelease() error {
	return fmt.Errorf("downloadRelease not implemented")
}

func (b *Borges) cacheName() string {
	binName := fmt.Sprintf("borges.%s", b.Version)
	return filepath.Join(b.binCache, binName)
}

func (b *Borges) tarName() string {
	return fmt.Sprintf("borges_%s_%s_amd64.tar.gz", b.Version, b.Distro)
}

func (b *Borges) dirName() string {
	return fmt.Sprintf("borges_%s_amd64", b.Version)
}

func fileExist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
