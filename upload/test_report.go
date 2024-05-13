package upload

import (
	"time"

	"github.com/codeclimate/test-reporter/formatters"
)

func NewTestReport(rep formatters.Report) *TestReport {
	tr := &TestReport{
		Type: "test_reports",
		Attributes: Attributes{
			RunAt:           time.Now().Unix(),
			CoveredPercent:  rep.CoveredPercent,
			CoveredStrength: rep.CoveredStrength,
			LineCounts:      rep.LineCounts,
		},
		SourceFiles: []SourceFile{},
	}
	for _, sf := range rep.SourceFiles {
		tr.SourceFiles = append(tr.SourceFiles, SourceFile{
			Type:            "test_file_reports",
			Coverage:        sf.Coverage,
			CoveredPercent:  sf.CoveredPercent,
			CoveredStrength: sf.CoveredStrength,
			LineCounts:      sf.LineCounts,
			Path:            sf.Name,
		})
	}
	return tr
}

type TestReport struct {
	Type        string       `json:"type"`
	Attributes  Attributes   `json:"attributes"`
	SourceFiles []SourceFile `json:"-"`
}
