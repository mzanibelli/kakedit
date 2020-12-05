package main

import (
	"fmt"
	"kakedit"
	"log"
	"os"
	"time"
)

const timeout time.Duration = 200 * time.Millisecond

func main() {
	var err error

	switch len(os.Args) {
	case 2: // Program invoked by the user
		err = kakedit.Pick(os.Args[0], os.Args[1],
			os.Getenv("kak_session"), os.Getenv("kak_client"), timeout)
	case 3: // Internal use only
		err = kakedit.Edit(os.Args[0], os.Args[1], os.Args[2])
	default:
		printUsage()
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: %s PROGRAM", os.Args[0])
	os.Exit(1)
}
