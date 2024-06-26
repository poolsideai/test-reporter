package dotcover

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	assert := require.New(t)

	formatter := Formatter{
		Path: "./example.xml",
	}
	rep, err := formatter.Format()
	assert.NoError(err)

	assert.Len(rep.SourceFiles, 3)
	assert.InDelta(71, rep.CoveredPercent, 1)
	assert.Equal(24, rep.LineCounts.Total)
	assert.Equal(17, rep.LineCounts.Covered)

	sf_one := rep.SourceFiles[`C:\Users\fulano\Desktop\unit-testing-using-mstest\PrimeService\PrimeService.cs`]
	assert.InDelta(83, sf_one.CoveredPercent, 1)

	sf_two := rep.SourceFiles[`C:\Users\fulano\Desktop\unit-testing-using-mstest\PrimeService\SecondService.cs`]
	assert.Equal(0.0, sf_two.CoveredPercent)

	sf_three := rep.SourceFiles[`C:\Users\fulano\Desktop\unit-testing-using-mstest\PrimeService.Tests\PrimeService_IsPrimeShould.cs`]
	assert.Equal(100.0, sf_three.CoveredPercent)
}
