package lcovjson

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Format(t *testing.T) {

	r := require.New(t)

	rb := Formatter{
		Path: "./lcovjson_example.json",
	}
	rep, err := rb.Format()
	r.NoError(err)

	r.InDelta(rep.CoveredPercent, 88.8, 1)

	sf := rep.SourceFiles["/Users/paulo/Development/GitHub/paulofaria/Codecov/Sources/Codecov/User.swift"]

	r.InDelta(sf.CoveredPercent, 66.66, 1)
	sfLc := sf.LineCounts
	r.Equal(sfLc.Covered, 6)
	r.Equal(sfLc.Missed, 3)
	r.Equal(sfLc.Total, 9)

	sf = rep.SourceFiles["/Users/paulo/Development/GitHub/paulofaria/Codecov/Tests/CodecovTests/CodecovTests.swift"]
	r.InDelta(sf.CoveredPercent, 100, 1)
	sfLc = sf.LineCounts
	r.Equal(sfLc.Covered, 18)
	r.Equal(sfLc.Missed, 0)
	r.Equal(sfLc.Total, 18)
}
