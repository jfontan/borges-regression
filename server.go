package regression

import (
	"fmt"
	"os/exec"
	"syscall"
)

type Server struct {
	cmd      *exec.Cmd
	repoPath string
	port     int
}

func NewServer(path string) (*Server, error) {
	return &Server{
		repoPath: path,
		port:     9418,
	}, nil
}

func (s *Server) Start() error {
	basePath := fmt.Sprintf("--base-path=%s", s.repoPath)
	port := fmt.Sprintf("--port=%d", s.port)

	s.cmd = exec.Command("git", "daemon", "--reuseaddr", basePath, port,
		"--export-all", s.repoPath)

	return s.cmd.Start()
}

func (s *Server) Stop() error {
	err := s.cmd.Process.Kill()
	if err != nil {
		return err
	}

	_ = s.cmd.Wait()
	return nil
}

func (s *Server) Alive() bool {
	if s.cmd == nil || s.cmd.Process == nil {
		return false
	}

	err := s.cmd.Process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}
