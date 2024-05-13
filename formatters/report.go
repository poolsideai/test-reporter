package formatters

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
)

type Report struct {
	CoveredPercent  float64     `json:"covered_percent"`
	CoveredStrength int         `json:"covered_strength"`
	LineCounts      LineCounts  `json:"line_counts"`
	SourceFiles     SourceFiles `json:"source_files"`
	RepoToken       string      `json:"repo_token"`
}

func NewReport() (Report, error) {
	rep := Report{
		SourceFiles: SourceFiles{},
		LineCounts:  LineCounts{},
	}

	return rep, nil
}

func (a *Report) Merge(reps ...*Report) error {
	for _, r := range reps {
		for _, sf := range r.SourceFiles {
			err := a.AddSourceFile(sf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (rep *Report) AddSourceFile(sf SourceFile) error {
	var err error

	// check if we already know about this file
	if s, ok := rep.SourceFiles[sf.Name]; ok {
		// remove the old values... we know more now
		rep.LineCounts.Covered -= s.LineCounts.Covered
		rep.LineCounts.Missed -= s.LineCounts.Missed
		rep.LineCounts.Total -= s.LineCounts.Total

		sf, err = s.Merge(sf)
		if err != nil {
			return err
		}
	} else {
		sf.CalcLineCounts()
	}
	rep.SourceFiles[sf.Name] = sf
	rep.LineCounts.Covered += sf.LineCounts.Covered
	rep.LineCounts.Missed += sf.LineCounts.Missed
	rep.LineCounts.Total += sf.LineCounts.Total

	rep.CoveredPercent = rep.LineCounts.CoveredPercent()
	return nil
}

func (r Report) Save(w io.Writer) error {
	b, err := json.MarshalIndent(r, "", "  ")
	logrus.Debugf("codeclimate.json content: %s", string(b))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
