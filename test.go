package regression

import "fmt"

type packResults map[string]*PackResult
type versionResults map[string]packResults

type Test struct {
	repos    *Repositories
	server   *Server
	versions []string
	borges   map[string]*Borges
	results  versionResults
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

var complexity = 1

func (t *Test) Run() error {
	results := make(versionResults)

	for _, version := range t.versions {
		_, ok := results[version]
		if !ok {
			results[version] = make(packResults)
		}

		borges, ok := t.borges[version]
		if !ok {
			panic("borges not initialized. Was Prepare called?")
		}

		fmt.Printf("## Version %s\n", version)

		for _, repo := range t.repos.Names(complexity) {
			url := t.server.Url(repo)
			pack, err := NewPack(borges.Path, url)
			if err != nil {
				return err
			}

			err = pack.Run()
			if err != nil {
				return err
			}

			results[version][repo] = pack.Result()

			fmt.Printf("  Repo: %s\n", repo)
			fmt.Printf("  Wall: %v\n", pack.wall)
			fmt.Printf("  Memory: %v\n", pack.rusage.Maxrss)
			fmt.Printf("  Files: %+v\n", pack.files)
		}
	}

	t.results = results

	return nil
}

func (t *Test) GetResults() bool {
	if len(t.versions) < 2 {
		panic("there should be at least two versions")
	}

	ok := true
	for i, version := range t.versions[0 : len(t.versions)-1] {
		fmt.Printf("#### Comparing %s - %s ####\n", version, t.versions[i+1])
		a := t.results[t.versions[i]]
		b := t.results[t.versions[i+1]]

		for _, repo := range t.repos.Names(complexity) {
			fmt.Printf("## Repo %s ##\n", repo)

			c := a[repo].ComparePrint(b[repo], 10.0)
			if !c {
				ok = false
			}
		}
	}

	return ok
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
