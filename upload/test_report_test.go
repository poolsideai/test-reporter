package upload

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func exampleReport() (formatters.Report, error) {
	rep := formatters.Report{}
	f, err := os.Open("../integration-tests/codeclimate.json")
	if err != nil {
		return rep, errors.WithStack(err)
	}

	rep, err = formatters.NewReport()
	if err != nil {
		return rep, errors.WithStack(err)
	}
	err = json.NewDecoder(f).Decode(&rep)
	if err != nil {
		return rep, errors.WithStack(err)
	}
	return rep, nil
}

func Test_NewTestReport(t *testing.T) {
	r := require.New(t)

	rep, err := exampleReport()
	r.NoError(err)
	data := NewTestReport(rep)
	r.Equal("test_reports", data.Type)

	at := data.Attributes
	r.NotZero(at.RunAt)
	r.InDelta(88.92, at.CoveredPercent, 1.0)
	r.Equal(0, at.CoveredStrength)
	r.Equal(rep.LineCounts, at.LineCounts)

	r.Len(data.SourceFiles, 20)
}
