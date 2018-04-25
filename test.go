package regression

import (
	"fmt"

	"gopkg.in/src-d/go-log.v0"
)

type packResults map[string][]*PackResult
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
var times = 3

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
		l, _ := log.New()
		l = l.New(log.Fields{"version": version})

		l.Debugf("Running version tests")

		for _, repo := range t.repos.Names(complexity) {
			results[version][repo] = make([]*PackResult, times)
			for i := 0; i < times; i++ {
				// TODO: do not stop on errors

				result, err := t.runTest(borges, repo)
				results[version][repo][i] = result

				if err != nil {
					return err
				}
			}
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

			// TODO: add more options like discard the first run, do the media, etc

			repoA, repoB := getResultsSmaller(a[repo], b[repo])

			c := repoA.ComparePrint(repoB, 10.0)
			if !c {
				ok = false
			}
		}
	}

	return ok
}

func (t *Test) runTest(borges *Borges, repo string) (*PackResult, error) {
	url := t.server.Url(repo)
	l, _ := log.New()
	log.Debugf("Executing pack test for %s", repo)

	pack, err := NewPack(borges.Path, url)
	if err != nil {
		log.Error(err, "Could not execute pack")
		return nil, err
	}

	err = pack.Run()
	out, _ := pack.Out()
	if err != nil {
		l.New(log.Fields{
			"repo":   repo,
			"borges": borges.Path,
			"url":    url,
			"output": out}).Error(err, "Could not execute pack")
		return nil, err
	}

	var fileSize int64
	for _, f := range pack.files {
		fileSize += f.Size()
	}

	l.New(log.Fields{
		"wall":     pack.wall,
		"memory":   pack.rusage.Maxrss,
		"fileSize": fileSize,
	}).Infof("finished pack")

	return pack.Result(), nil
}

func (t *Test) prepareServer() error {
	log.Infof("Downloading repositories")
	err := t.repos.Download()
	if err != nil {
		return err
	}

	log.Infof("Starting git server")
	err = t.server.Start()
	return err
}

func (t *Test) prepareBorges() error {
	log.Infof("Preparing borges binaries")
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

// Get the runs with lower wall time
func getResultsSmaller(
	a []*PackResult,
	b []*PackResult,
) (*PackResult, *PackResult) {
	repoA := a[0]
	repoB := b[0]
	for i := 1; i < len(a); i++ {
		if a[i].Wtime < repoA.Wtime {
			repoA = a[i]
		}

		if b[i].Wtime < repoB.Wtime {
			repoB = b[i]
		}
	}

	return repoA, repoB
}
