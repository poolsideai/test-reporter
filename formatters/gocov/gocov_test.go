package gocov

import (
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	t.Run("should parse report from single package", func(t *testing.T) {

		r := require.New(t)

		f := &Formatter{Path: "./example.out"}
		rep, err := f.Format()
		r.NoError(err)

		r.Len(rep.SourceFiles, 4)

		sf, fok := rep.SourceFiles[filepath.FromSlash("github.com/codeclimate/test-reporter/formatters/source_file.go")]
		r.True(fok)

		r.InDelta(75.8, sf.CoveredPercent, 1)
		r.Len(sf.Coverage, 115)
		r.False(sf.Coverage[5].Valid)
		r.True(sf.Coverage[54].Valid)
		r.Equal(0, sf.Coverage[52].Int)
		r.Equal(1, sf.Coverage[55].Int)
	})

	//
	// Test parsing report that was generated using `-coverpkg` flag.
	//
	// go test \
	// -coverpkg="github.com/codeclimate/test-reporter/formatters/gocov/example/foo,github.com/codeclimate/test-reporter/formatters/gocov/example/bar" \
	// -coverprofile=example_foobar.out \
	// ./...
	//
	t.Run("should parse coverage report from multiple packages", func(t *testing.T) {

		r := require.New(t)

		f := &Formatter{Path: filepath.Join("example", "foobar_test.out"), GoModuleName: "github.com/codeclimate/test-reporter"}
		rep, err := f.Format()
		r.NoError(err)

		r.Len(rep.SourceFiles, 2)

		sfFoo, okFoo := rep.SourceFiles[filepath.Join("formatters", "gocov", "example", "foo", "foo.go")]
		r.True(okFoo)
		sfBar, okBar := rep.SourceFiles[filepath.Join("formatters", "gocov", "example", "bar", "bar.go")]
		r.True(okBar)

		r.EqualValues(85, rep.CoveredPercent)
		r.EqualValues(100, sfFoo.CoveredPercent)
		r.InDelta(66.66, sfBar.CoveredPercent, 0.01)
	})
}
