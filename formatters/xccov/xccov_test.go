package xccov

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Format(t *testing.T) {

	r := require.New(t)

	rb := Formatter{
		Path: "./xccov_example.json",
	}
	rep, err := rb.Format()
	r.NoError(err)

	r.InDelta(rep.CoveredPercent, 20.93, 1)

	sf := rep.SourceFiles["Documents/github/ww/ios-SuperApp/Pods/BundleA/BundleA/BundleA.m"]
	r.Equal(sf.CoveredPercent, 0.0)
	sfLc := sf.LineCounts
	r.Equal(sfLc.Covered, 0)
	r.Equal(sfLc.Missed, 10)
	r.Equal(sfLc.Total, 10)

	sf = rep.SourceFiles["Documents/github/ww/ios-SuperApp/Pods/SuperClass/SuperClass/SuperClass.m"]
	r.InDelta(sf.CoveredPercent, 22.68, 1)
	sfLc = sf.LineCounts
	r.Equal(sfLc.Covered, 27)
	r.Equal(sfLc.Missed, 92)
	r.Equal(sfLc.Total, 119)
}
