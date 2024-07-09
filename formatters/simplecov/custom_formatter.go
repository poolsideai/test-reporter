package simplecov

import (
	"encoding/json"
	"github.com/codeclimate/test-reporter/formatters"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func customFormat(r Formatter, rep formatters.Report) (formatters.Report, error) {
	logrus.Debugf("Analyzing simplecov json output from custom format %s", r.Path)
	jf, err := os.Open(r.Path)
	if err != nil {
		return rep, errors.WithStack(errors.Errorf("could not open coverage file %s", r.Path))
	}

	m := map[string]simplecovJsonFormatterReport{}
	err = json.NewDecoder(jf).Decode(&m)
	if err != nil {
		return rep, err
	}

	for _, v := range m {
		for n, ls := range v.CoverageType {
			fe, err := formatters.NewSourceFile(n, nil)
			if err != nil {
				return rep, errors.WithStack(err)
			}
			fe.Coverage = transformLineCoverageToCoverage(ls.LineCoverage)
			err = rep.AddSourceFile(fe)
			if err != nil {
				return rep, errors.WithStack(err)
			}
		}
	}

	return rep, nil
}
