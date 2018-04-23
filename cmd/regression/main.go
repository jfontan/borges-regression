package main

import (
	"fmt"

	"github.com/jfontan/borges-regression"

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

	test, err := regression.NewTest(args)
	if err != nil {
		panic(err)
	}

	err = test.Prepare()
	if err != nil {
		panic(err)
	}

	err = test.Run()
	if err != nil {
		panic(err)
	}

	err = test.Stop()
	if err != nil {
		panic(err)
	}
}
