package regression

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecution(t *testing.T) {
	require := require.New(t)

	e, err := NewExecutor("echo", "test", "message")
	require.NoError(err)

	_, err = e.Out()
	require.Equal(ErrNotRun, err)

	err = e.Run()
	require.NoError(err)

	out, err := e.Out()
	require.NoError(err)
	require.Equal("test message\n", out)

	rusage, err := e.Rusage()
	require.NoError(err)

	wall, err := e.Wall()
	require.NoError(err)

	require.True(rusage.Maxrss > 0)
	require.True(rusage.Utime.Nano() > 0)
	require.True(wall > 0*time.Second)
}
