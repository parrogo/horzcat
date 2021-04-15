package horzcat

import (
	"bytes"
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures
var fixtureRootFS embed.FS
var fixtureFS, _ = fs.Sub(fixtureRootFS, "fixtures")

func TestPlaceholder(t *testing.T) {
	content, err := fs.ReadFile(fixtureFS, "placeholder")
	if assert.NoError(t, err) {
		assert.Equal(t, "this is a placeholder", string(content))
	}
}

func TestConcat(t *testing.T) {
	//assert := assert.New(t)

	t.Run("Fail if no sources are provided", func(t *testing.T) {
		require := require.New(t)

		err := Concat(Options{})
		require.EqualError(err, "no source readers provided")
	})

	t.Run("Append opt.Tail with 1 source reader alone", func(t *testing.T) {
		require := require.New(t)

		source, err := fixtureFS.Open("lines1.txt")
		require.NoError(err)
		var buf bytes.Buffer
		err = Concat(Options{
			Tail:   "!!",
			Target: &buf,
		}, source)
		require.NoError(err)

		actual := buf.String()
		expected := "ciao!!\nsalve!!\n"

		require.Equal(expected, actual)
	})

}
