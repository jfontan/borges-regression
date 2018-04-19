package regression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecution(t *testing.T) {
	require := require.New(t)

	e, err := NewExecutor("echo", "hola", "caracola")
	require.NoError(err)

	_, err = e.Out()
	require.Equal(ErrNotRun, err)

	err = e.Run()
	require.NoError(err)

	out, err := e.Out()
	require.NoError(err)
	require.Equal("hola caracola\n", out)

	rusage, err := e.Rusage()
	require.NoError(err)

	wall, err := e.Wall()
	require.NoError(err)

	fmt.Printf("stime: %v\n", rusage.Stime)
	fmt.Printf("utime: %v\n", rusage.Utime)
	fmt.Printf("wall: %v\n", wall)
}
