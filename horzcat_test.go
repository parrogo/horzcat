package horzcat

import (
	"bytes"
	"embed"
	"errors"
	"io"
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

type FailingReader struct {
}

func (w FailingReader) Read(but []byte) (int, error) {
	return 0, errors.New("expected error")
}

var _ io.Reader = FailingReader{}

func TestConcat(t *testing.T) {
	//assert := assert.New(t)

	t.Run("Fail if no sources are provided", func(t *testing.T) {
		require := require.New(t)

		err := Concat(Options{})
		require.EqualError(err, "no source readers provided")
	})

	t.Run("InputError can be converted to include filenames", func(t *testing.T) {
		origErr := errors.New("expected error")
		err := InputError{origErr, 1}
		msg := err.Convert([]string{"file1", "file2", "file3"})
		assert.Equal(t, "Cannot read file file2: expected error", msg.Error())
	})

	t.Run("Return InputError on failing readers", func(t *testing.T) {
		require := require.New(t)
		source1, err := fixtureFS.Open("lines1.txt")
		require.NoError(err)

		err = Concat(Options{
			Target: io.Discard,
		}, source1, FailingReader{})
		require.EqualError(err, "error reading from source 1: expected error")
		var idxErr InputError
		require.ErrorAs(err, &idxErr)

		require.Equal(1, idxErr.idx)
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

	t.Run("Skip column headers", func(t *testing.T) {
		require := require.New(t)

		source1, err := fixtureFS.Open("lines1.txt")
		require.NoError(err)
		source2, err := fixtureFS.Open("lines3.txt")
		require.NoError(err)

		var buf bytes.Buffer
		err = Concat(Options{
			RowHeaderLen: 2,
			Target:       &buf,
		}, source1, source2)
		require.NoError(err)

		actual := buf.String()
		expected := "ciaodre\nsalveP\nThe end\n"

		require.Equal(expected, actual)
	})

	t.Run("Append opt.Tail & opt.Sep with multiple source readers", func(t *testing.T) {
		require := require.New(t)

		source1, err := fixtureFS.Open("lines1.txt")
		require.NoError(err)
		source2, err := fixtureFS.Open("lines2.txt")
		require.NoError(err)

		var buf bytes.Buffer
		err = Concat(Options{
			Tail:   "😎",
			Sep:    " ",
			Target: &buf,
		}, source1, source2)
		require.NoError(err)

		actual := buf.String()
		expected := "ciao Andre😎\nsalve Parro😎\nThe end😎\n"

		require.Equal(expected, actual)
	})

}
