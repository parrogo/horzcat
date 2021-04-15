package horzcat_test

import (
	"embed"
	"fmt"
	"io/fs"
	
	"github.com/parrogo/horzcat"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

// This example show how to use horzcat.Func()
func ExampleFunc() {
	fmt.Println(horzcat.Func())
	// Output: 42
}
