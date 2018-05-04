package regression

import (
	"os/exec"
	"syscall"
)

type Server struct {
	cmd *exec.Cmd
}

func NewServer() *Server {
	return new(Server)
}

func (s *Server) Start(name string, arg ...string) error {
	s.cmd = exec.Command(name, arg...)

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
