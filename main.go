package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kdisneur/yamlsort/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var flagCfg struct {
		Indent         int
		DisplayVersion bool
	}

	flag.CommandLine.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: yamlsort [-indent 2]")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "DESCRIPTION")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "  read a YAML or partial YAML from STDIN, sort the keys and send them back to STDOUT")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "FLAGS")
		flag.CommandLine.PrintDefaults()

	}

	flag.IntVar(&flagCfg.Indent, "indent", 2, "number of spaces to use for indentation")
	flag.BoolVar(&flagCfg.DisplayVersion, "v", false, "display the command line version")
	flag.Parse()

	if flagCfg.DisplayVersion {
		fmt.Println(internal.GetVersionInfo())
		return nil
	}

	return internal.SortYAML(os.Stdin, os.Stdout, flagCfg.Indent)
}
