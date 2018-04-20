package regression

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getGopath() (string, error) {
	gopath := os.Getenv("GOPATH")
	split := strings.Split(gopath, ":")

	if len(split) > 0 {
		return split[0], nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func TestServer(t *testing.T) {
	require := require.New(t)

	gopath, err := getGopath()
	require.NoError(err)

	repos := fmt.Sprintf("%s/src/github.com/src-d", gopath)
	server, err := NewServer(repos)
	require.NoError(err)

	err = server.Start()
	require.NoError(err)
	require.True(server.Alive())

	err = server.Stop()
	require.NoError(err)
	require.False(server.Alive())
}
