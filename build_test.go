package regression

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	require := require.New(t)

	version := "remote:master"

	build, err := NewBuild(NewConfig(), version)
	require.NoError(err)

	_, err = build.download()
	require.NoError(err)

	err = build.build()
	require.NoError(err)
}
