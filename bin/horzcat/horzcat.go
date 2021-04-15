package main

import (
	"flag"
	"fmt"
	"os"
)

// Version of the command
var Version string = "development"

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	os.Exit(1)
}

var options struct {
	version bool
	sep     string
	tail    string
}

func usage(msg string) {
	fmt.Fprintf(os.Stderr, "Wrong usage: %s\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.BoolVar(&options.version, "v", false, "print version of the command to stdout")

	flag.StringVar(&options.sep, "s", "", "separator added between lines")
	flag.StringVar(&options.tail, "t", "", "tail string add at end of every concateneted line")

	flag.Parse()

	if options.version {
		fmt.Println(Version)
		os.Exit(0)
	}
}
