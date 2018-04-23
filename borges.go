package regression

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alcortesm/tgz"
	"gopkg.in/src-d/go-errors.v1"
)

var regRelease = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

var ErrBinaryNotFound = errors.NewKind(
	"borges binary not found in release tarball")

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
	tmpDir, err := createTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	r := GetReleases()

	download := filepath.Join(tmpDir, "download.tar.gz")
	err = r.Get(b.Version, b.tarName(), download)
	if err != nil {
		return err
	}

	path, err := tgz.Extract(download)
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)

	binary := filepath.Join(path, b.dirName(), "borges")
	exist, err := fileExist(binary)
	if err != nil {
		return err
	}
	if !exist {
		return ErrBinaryNotFound.New()
	}

	orig, err := os.Open(binary)
	if err != nil {
		return err
	}

	os.MkdirAll(b.binCache, 0755)
	dst, err := os.Create(b.cacheName())
	if err != nil {
		return err
	}
	dst.Chmod(0755)
	defer dst.Close()

	_, err = io.Copy(dst, orig)
	if err != nil {
		dst.Close()
		os.Remove(dst.Name())
		return err
	}

	return nil
}

func (b *Borges) cacheName() string {
	binName := fmt.Sprintf("borges.%s", b.Version)
	return filepath.Join(b.binCache, binName)
}

func (b *Borges) tarName() string {
	return fmt.Sprintf("borges_%s_%s_amd64.tar.gz", b.Version, b.Distro)
}

func (b *Borges) dirName() string {
	return fmt.Sprintf("borges_%s_amd64", b.Distro)
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
