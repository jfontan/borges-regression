package regression

import "fmt"

type Test struct {
	repos    *Repositories
	server   *Server
	versions []string
	borges   map[string]*Borges
}

func NewTest(versions []string) (*Test, error) {
	repos := NewRepositories()
	server, err := NewServer(repos.Path())
	if err != nil {
		return nil, err
	}

	return &Test{
		repos:    repos,
		server:   server,
		versions: versions,
	}, nil
}

func (t *Test) Prepare() error {
	err := t.prepareServer()
	if err != nil {
		return err
	}

	err = t.prepareBorges()
	return err
}

func (t *Test) Stop() error {
	return t.server.Stop()
}

func (t *Test) Run() error {
	for _, version := range t.versions {
		borges, ok := t.borges[version]
		if !ok {
			panic("borges not initialized. Was Prepare called?")
		}

		fmt.Printf("## Version %s\n", version)

		for _, repo := range t.repos.Names(1) {
			url := t.server.Url(repo)
			pack, err := NewPack(borges.Path, url)
			if err != nil {
				return err
			}

			err = pack.Run()
			if err != nil {
				return err
			}

			fmt.Printf("  Repo: %s\n", repo)
			fmt.Printf("  Wall: %v\n", pack.wall)
			fmt.Printf("  Memory: %v\n", pack.rusage.Maxrss)
			fmt.Printf("  Files: %+v\n", pack.files)
		}
	}

	return nil
}

func (t *Test) prepareServer() error {
	err := t.repos.Download()
	if err != nil {
		return err
	}

	err = t.server.Start()
	return err
}

func (t *Test) prepareBorges() error {
	t.borges = make(map[string]*Borges, len(t.versions))
	for _, version := range t.versions {
		b := NewBorges(version)
		err := b.Download()
		if err != nil {
			return err
		}

		t.borges[version] = b
	}

	return nil
}
