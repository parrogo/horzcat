// Package horzcat concatenates multiple files in horizontal direction.
// To clarify, consider the following comparison between cat and horzcat of two files:
//   $ cat a.dat 1.dat
//   a b
//   c d
//   1 2
//   3 4
//   $ horzcat -s ' ' a.dat 1.dat
//   a b 1 2
//   c d 3 4
package horzcat

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// Options struct groups all options
// accepted by Concat.
//
// Target field contains the io.Writer
// on which to write concateneted lines.
// When it's nil, os.Stdout is used as writer.
//
// Sep field is a string used to separate
// lines from source readers.
//
// Tail field is a string appended
// at the end of every concatenated line.
//
// SameLinesCount field, when set
// to true, requires that every source
// has the exact numer of lines.
// If one or more of the readers has a different
// lines count, an error is returned.
// When the field is false, excess lines from
// one or more reader are still concatened and written
// to output. Sep string is added alone for each of the
// readers that miss one or more lines.
type Options struct {
	Target         io.Writer
	Sep            string
	Tail           string
	SameLinesCount bool
}

// Concat read lines from all io.Reader in sources,
// concatenetes line by line horizontally, and finally writes
// concateneted lines to options.Target argument.
// If options.Target is nil, Concat writes to os.Stdout.
// String opt.Sep is added between lines
// String opt.Tail is added at the end of every written line.
func Concat(opt Options, sources ...io.Reader) error {
	if len(sources) == 0 {
		return errors.New("no source readers provided")
	}
	var out *bufio.Writer
	if opt.Target != nil {
		out = bufio.NewWriter(opt.Target)
	} else {
		out = bufio.NewWriter(os.Stdout)
	}

	bufreaders := make([]*bufio.Scanner, len(sources))
	for idx, source := range sources {
		bufreaders[idx] = bufio.NewScanner(source)
	}

	source := bufreaders[0]

	for {
		sourceok := source.Scan()
		if !sourceok {
			break
		}

		_, err := out.Write(source.Bytes())
		if err != nil {
			return err
		}

		_, err = out.WriteString(opt.Tail + "\n")
		if err != nil {
			return err
		}
	}

	err := out.Flush()
	if err != nil {
		return err
	}

	return nil
}
