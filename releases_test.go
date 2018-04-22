package regression

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleases(t *testing.T) {
	require := require.New(t)
	dir, err := createTempDir()
	require.NoError(err)
	defer os.RemoveAll(dir)

	require.Nil(releases)
	r := GetReleases()
	require.NotNil(releases)
	require.Equal(releases, r)
	require.Nil(r.repoReleases)

	path := filepath.Join(dir, "invalid_version")
	err = r.Get("invalid_version", "invalid_asset", path)
	require.Error(err)
	require.True(ErrVersionNotFound.Is(err))
	require.False(fileExist(path))

	list := r.repoReleases

	path = filepath.Join(dir, "invalid_asset")
	err = r.Get("v0.12.0", "invalid_asset", path)
	require.Error(err)
	require.True(ErrAssetNotFound.Is(err))
	require.False(fileExist(path))

	require.Exactly(list, r.repoReleases)

	path = filepath.Join(dir, "borges.v0.12.0")
	err = r.Get("v0.12.0", "borges_v0.12.0_linux_amd64.tar.gz", path)
	require.NoError(err)
	require.True(fileExist(path))
}
