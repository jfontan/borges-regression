package main

import (
	"fmt"

	"github.com/jfontan/borges-regression"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/jessevdk/go-flags.v1"
)

type packCmd struct {
}

func main() {
	parser := flags.NewParser(nil, flags.Default)
	parser.LongDescription = "long description"

	args, err := parser.Parse()
	fmt.Printf("Args: %+v\n", args)
	fmt.Printf("Error: %+v\n", err)

	repos := regression.NewRepositories()
	err = repos.Download()
	if err != nil {
		panic(err)
	}

	var borges []*regression.Borges

	for _, v := range args {
		fmt.Printf("Processing %s\n", v)
		b := regression.NewBorges(v)
		err := b.Download()
		if err != nil {
			panic(err)
		}

		borges = append(borges, b)
	}

	spew.Dump(borges)
}
