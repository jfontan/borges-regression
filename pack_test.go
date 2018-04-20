package regression

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPack(t *testing.T) {
	require := require.New(t)

	p, err := NewPack("borges", "https://github.com/jfontan/cangallo.git")
	require.NoError(err)

	err = p.Run()
	require.NoError(err)

	output, err := p.Out()
	require.NoError(err)
	require.NotEmpty(output)

	require.True(p.test)
	require.Len(p.files, 1)
}
