package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/parrogo/horzcat"
)

// Version of the command
var Version string = "development"

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

var options struct {
	version bool
	sep     string
	tail    string
	out     string
}

func usage(msg string) {
	fmt.Fprintf(os.Stderr, "Wrong usage: %s\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {
	flag.BoolVar(&options.version, "v", false, "print version of the command to stdout.")

	flag.StringVar(&options.sep, "s", "", "separator added between lines.")
	flag.StringVar(&options.tail, "t", "", "tail string add at end of every concateneted line.")
	flag.StringVar(&options.out, "out", "", "name of output file. Defaults to stdout.")

	flag.Parse()

	if options.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	sources := make([]io.Reader, len(flag.Args()))
	for idx, arg := range flag.Args() {
		f, err := os.Open(arg)
		fatal(err)
		sources[idx] = f
		defer func(f io.Closer) {
			fatal(f.Close())
		}(f)
	}

	opt := horzcat.Options{
		Sep:  options.sep,
		Tail: options.tail,
	}

	if options.out != "" {
		f, err := os.OpenFile(options.out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
		fatal(err)
		opt.Target = f
		defer func(f io.Closer) {
			fatal(f.Close())
		}(f)
	}

	err := horzcat.Concat(opt, sources...)
	fatal(err)
}
