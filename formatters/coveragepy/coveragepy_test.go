package coveragepy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {

	r := require.New(t)

	f := &Formatter{Path: "./example.xml"}
	rep, err := f.Format()
	r.NoError(err)
	r.Len(rep.SourceFiles, 12)

	sf := rep.SourceFiles["codeclimate_test_reporter/components/runner.py"]
	r.InDelta(85.71, sf.CoveredPercent, 1)
	r.Len(sf.Coverage, 82)
	r.False(sf.Coverage[53].Valid)
	r.True(sf.Coverage[54].Valid)
	r.Equal(1, sf.Coverage[54].Int)
	r.Equal(0, sf.Coverage[55].Int)
}
