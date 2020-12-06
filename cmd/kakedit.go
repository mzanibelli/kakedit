package main

import (
	"flag"
	"fmt"
	"kakedit"
	"log"
	"os"
)

var local bool

func init() {
	flag.BoolVar(&local, "local", false,
		"run a new instance instead of sending a command to an existing one")
}

func main() {
	var err error

	flag.Parse()

	switch {
	case len(flag.Args()) == 1 && local:
		err = kakedit.Local(flag.Arg(0))
	case len(flag.Args()) == 1:
		err = kakedit.Server(flag.Arg(0))
	case len(flag.Args()) == 2: // Internal use only
		err = kakedit.Client(flag.Arg(0), flag.Arg(1))
	default:
		printUsage()
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-local] PROGRAM", os.Args[0])
	os.Exit(1)
}
