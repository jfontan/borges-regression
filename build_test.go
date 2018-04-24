package regression

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	require := require.New(t)

	version := "local:master"

	build, err := NewBuild(version)
	require.NoError(err)

	err = build.download()
	require.NoError(err)
}
