package regression

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alcortesm/tgz"
	errors "gopkg.in/src-d/go-errors.v0"
	log "gopkg.in/src-d/go-log.v0"
)

var regRelease = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// ErrBinaryNotFound is returned when the executable is not found in
// the release tarball.
var ErrBinaryNotFound = errors.NewKind("binary not found in release tarball")

// Binary struct contains information and functionality to prepare and
// use a binary version.
type Binary struct {
	Name    string
	Version string
	Path    string

	releases *Releases
	config   Config
}

// NewBinary creates a new Binary structure.
func NewBinary(
	config Config,
	name, version string,
	releases *Releases,
) *Binary {
	return &Binary{
		Name:     name,
		Version:  version,
		releases: releases,
		config:   config,
	}
}

// IsRelease checks if the version matches the format of a release, for
// example v0.12.1.
func (b *Binary) IsRelease() bool {
	return regRelease.MatchString(b.Version)
}

// Download prepares a binary version if it's still not in the
// binaries directory.
func (b *Binary) Download() error {
	switch {
	case IsRepo(b.Version):
		build, err := NewBuild(b.config, b.Version)
		if err != nil {
			return err
		}

		binary, err := build.Build()
		if err != nil {
			return err
		}

		b.Path = binary
		return nil

	case b.Version == "latest":
		version, err := b.releases.Latest()
		if err != nil {
			return nil
		}

		b.Version = version

	case !b.IsRelease():
		b.Path = b.Version
		return nil
	}

	cacheName := b.cacheName()
	exist, err := fileExist(cacheName)
	if err != nil {
		return err
	}

	if exist {
		log.Debugf("Binary for %s already downloaded", b.Version)
		b.Path = cacheName
		return nil
	}

	log.Debugf("Dowloading version %s", b.Version)
	err = b.downloadRelease()
	if err != nil {
		log.Error(err, "Could not download version %s", b.Version)
		return err
	}

	b.Path = cacheName

	return nil
}

func (b *Binary) downloadRelease() error {
	tmpDir, err := createTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	download := filepath.Join(tmpDir, "download.tar.gz")
	err = b.releases.Get(b.Version, b.tarName(), download)
	if err != nil {
		return err
	}

	path, err := tgz.Extract(download)
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)

	binary := filepath.Join(path, b.dirName(), b.Name)
	err = copyBinary(binary, b.cacheName())

	return err
}

func (b *Binary) cacheName() string {
	binName := fmt.Sprintf("%s.%s", b.Name, b.Version)
	return filepath.Join(b.config.BinaryCache, binName)
}

func (b *Binary) tarName() string {
	return fmt.Sprintf("%s_%s_%s_amd64.tar.gz", b.Name, b.Version, b.config.OS)
}

func (b *Binary) dirName() string {
	return fmt.Sprintf("%s_%s_amd64", b.Name, b.config.OS)
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

func copyBinary(source, destination string) error {
	println("copyBinary", source, destination)
	exist, err := fileExist(source)
	if err != nil {
		return err
	}
	if !exist {
		return ErrBinaryNotFound.New()
	}

	orig, err := os.Open(source)
	if err != nil {
		return err
	}

	dir := filepath.Dir(destination)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	dst, err := os.Create(destination)
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
