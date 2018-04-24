package main

import (
	"os"

	"github.com/jfontan/borges-regression"

	flags "github.com/jessevdk/go-flags"
	"gopkg.in/src-d/go-log.v0"
)

type packCmd struct {
}

func main() {
	parser := flags.NewParser(nil, flags.Default)
	parser.LongDescription = "long description"

	args, err := parser.Parse()
	if err != nil {
		log.Error(err, "Could not parse arguments %v")
		os.Exit(1)
	}

	if len(args) < 2 {
		log.Error(nil, "There should be at least two versions")
		os.Exit(1)
	}

	test, err := regression.NewTest(args)
	if err != nil {
		panic(err)
	}

	log.Infof("Preparing run")
	err = test.Prepare()
	if err != nil {
		log.Error(err, "Could not prepare environment")
		os.Exit(1)
	}

	err = test.Run()
	if err != nil {
		panic(err)
	}

	res := test.GetResults()

	err = test.Stop()
	if err != nil {
		panic(err)
	}

	if !res {
		os.Exit(1)
	}
}
