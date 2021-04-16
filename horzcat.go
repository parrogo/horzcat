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
	"fmt"
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
		lineScan := bufio.NewScanner(source)
		buf := make([]byte, 0, 10*1024*1024)
		lineScan.Buffer(buf, 10*1024*1024)
		bufreaders[idx] = lineScan
	}

	lines := make([][]byte, len(sources))
	empty := make([]byte, 0)
	allSourcesDone := false

	outWriteErr := func(err error) error {
		return fmt.Errorf("cannot write to output: %w", err)
	}

	for !allSourcesDone {
		allSourcesDone = true
		for idx, source := range bufreaders {
			if source.Scan() {
				allSourcesDone = false
				lines[idx] = source.Bytes()
			} else {
				lines[idx] = empty
			}
		}

		if allSourcesDone {
			continue
		}

		somethingWritten := false
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			if somethingWritten {
				_, err := out.WriteString(opt.Sep)
				if err != nil {
					return outWriteErr(err)
				}
			}
			_, err := out.Write(line)
			if err != nil {
				return outWriteErr(err)
			}
			somethingWritten = true
		}

		_, err := out.WriteString(opt.Tail + "\n")
		if err != nil {
			return outWriteErr(err)
		}
	}

	err := out.Flush()
	if err != nil {
		return fmt.Errorf("cannot flush output: %w", err)
	}

	for idx, source := range bufreaders {
		if err := source.Err(); err != nil {
			return fmt.Errorf("error reading from source %d: %w", idx, InputError{err, idx})
		}
	}

	return nil
}

// InputError wraps an error
// in order to include the position
// of failing stream.
type InputError struct {
	err error
	idx int
}

// Error implements error interface
func (e InputError) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error
func (e InputError) Unwrap() error {
	return e.err
}

// Convert returns an error that include the
// name of the file that causes the failure
func (e InputError) Convert(filenames []string) error {
	return fmt.Errorf("Cannot read file %s: %w", filenames[e.idx], e.err)
}
