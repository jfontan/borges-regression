package regression

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Pack struct {
	*Executor
	test   bool
	binary string
	repo   string
	files  []os.FileInfo
}

func NewPack(binary, repo string) (*Pack, error) {
	return &Pack{
		Executor: new(Executor),
		binary:   binary,
		repo:     repo,
	}, nil
}

func (p *Pack) Run() error {
	list, err := createList(p.repo)
	if err != nil {
		return err
	}
	defer os.Remove(list)

	dir, err := createTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	lArg := fmt.Sprintf("--file=%s", list)
	dArg := fmt.Sprintf("--to=%s", dir)

	executor, err := NewExecutor(p.binary, "pack", lArg, dArg)
	if err != nil {
		return err
	}

	p.Executor = executor

	err = p.Executor.Run()
	if err != nil {
		return err
	}

	files, err := fileInfo(dir)
	if err != nil {
		return err
	}

	p.files = files
	p.test = true

	return nil
}

func (p *Pack) Files() ([]os.FileInfo, error) {
	if !p.executed {
		return nil, ErrNotRun
	}

	return p.files, nil
}

func createList(repo string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "packer-list")
	if err != nil {
		return "", err
	}

	_, err = tmpFile.WriteString(repo)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}

	err = tmpFile.Close()
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func createTempDir() (string, error) {
	dir, err := ioutil.TempDir("", "packer-dir")
	if err != nil {
		return "", err
	}

	return dir, nil
}

func fileInfo(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}
