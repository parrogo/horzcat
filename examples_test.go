package horzcat_test

import (
	"bytes"
	"embed"
	"io/fs"

	"github.com/parrogo/horzcat"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

// This example show how to use horzcat.Func()
func ExampleConcat() {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer

	buf1.WriteString("ciao\n")
	buf1.WriteString("salve\n")

	buf2.WriteString("Andre\n")
	buf2.WriteString("Parro\n")
	buf2.WriteString("The end\n")

	err := horzcat.Concat(horzcat.Options{
		Sep:  ",",
		Tail: "!",
	})
	if err != nil {
		panic(err)
	}
	// Output: ciao,Andre!
	// salve,Parro!
	// The end!
}
