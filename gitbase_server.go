package regression

type GitbaseServer struct {
	*Server
	binary string
	repos  string
}

func NewGitbaseServer(binary, repos string) *GitbaseServer {
	return &GitbaseServer{
		Server: NewServer(),
		binary: binary,
		repos:  repos,
	}
}

func (s *GitbaseServer) Start() error {
	return s.Server.Start(s.binary, "-g", s.repos)
}
