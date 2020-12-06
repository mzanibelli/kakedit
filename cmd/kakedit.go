package main

import (
	"flag"
	"kakedit"
	"kakedit/cli"
	"strings"
)

func main() {
	var err error

	if err = cli.Parse(); err != nil {
		cli.Fail(err)
	}

	switch cli.Mode {
	case cli.ModeLocal:
		err = kakedit.Local(strings.Join(flag.Args(), " "))
	case cli.ModeServer:
		err = kakedit.Server(strings.Join(flag.Args(), " "))
	case cli.ModeClient:
		err = kakedit.Client(flag.Arg(0), flag.Arg(1))
	default:
		cli.Fail(nil)
	}

	if err != nil {
		cli.Fail(err)
	}
}
