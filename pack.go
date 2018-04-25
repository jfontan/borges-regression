package regression

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
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

func (p *Pack) Result() *PackResult {
	var size int64

	for _, file := range p.files {
		size += file.Size()
	}

	packResult := &PackResult{
		Memory:   p.rusage.Maxrss,
		Wtime:    p.wall,
		Stime:    time.Duration(p.rusage.Stime.Nano()),
		Utime:    time.Duration(p.rusage.Utime.Nano()),
		Files:    p.files,
		FileSize: size,
	}

	return packResult
}

type PackResult struct {
	Memory   int64
	Wtime    time.Duration
	Stime    time.Duration
	Utime    time.Duration
	Files    []os.FileInfo
	FileSize int64
}

type PackComparison struct {
	Memory   float64
	Wtime    float64
	Stime    float64
	Utime    float64
	FileSize float64
}

func (p *PackResult) Compare(q *PackResult) PackComparison {
	return PackComparison{
		Memory:   percent(p.Memory, q.Memory),
		Wtime:    percent(int64(p.Wtime), int64(q.Wtime)),
		Stime:    percent(int64(p.Stime), int64(q.Stime)),
		Utime:    percent(int64(p.Utime), int64(q.Utime)),
		FileSize: percent(p.FileSize, q.FileSize),
	}
}

var compareFormat = "%s: %v -> %v (%v), %v\n"

func (p *PackResult) ComparePrint(q *PackResult, allowance float64) bool {
	ok := true
	c := p.Compare(q)

	if c.Memory > allowance {
		ok = false
	}
	fmt.Printf(compareFormat,
		"Memory",
		p.Memory,
		q.Memory,
		c.Memory,
		allowance > c.Memory,
	)

	if c.Wtime > allowance {
		ok = false
	}
	fmt.Printf(compareFormat,
		"Wtime",
		p.Wtime,
		q.Wtime,
		c.Wtime,
		allowance > c.Wtime,
	)

	// if c.Stime > allowance {
	// 	ok = false
	// }
	fmt.Printf(compareFormat,
		"Stime",
		p.Stime,
		q.Stime,
		c.Stime,
		allowance > c.Stime,
	)

	// if c.Utime > allowance {
	// 	ok = false
	// }
	fmt.Printf(compareFormat,
		"Utime",
		p.Utime,
		q.Utime,
		c.Utime,
		allowance > c.Utime,
	)

	if c.FileSize > allowance {
		ok = false
	}
	fmt.Printf(compareFormat,
		"FileSize",
		p.FileSize,
		q.FileSize,
		c.FileSize,
		allowance > c.FileSize,
	)

	return ok
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

func percent(a, b int64) float64 {
	diff := b - a
	return (float64(diff) / float64(a)) * 100
}
