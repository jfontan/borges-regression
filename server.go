package regression

import (
	"fmt"
	"os/exec"
	"syscall"
)

type Server struct {
	cmd    *exec.Cmd
	config Config
}

func NewServer(config Config) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}

func (s *Server) Start() error {
	basePath := fmt.Sprintf("--base-path=%s", s.config.RepositoriesCache)
	port := fmt.Sprintf("--port=%d", s.config.GitServerPort)

	s.cmd = exec.Command("git", "daemon", "--reuseaddr", basePath, port,
		"--export-all", s.config.RepositoriesCache)

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
	return err == nil
}

func (s *Server) Url(name string) string {
	return fmt.Sprintf("git://localhost:%d/%s", s.config.GitServerPort, name)
}
