package gcov

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codeclimate/test-reporter/formatters"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Formatter collects GCov files, parses them, then formats them into a single report.
type Formatter struct {
	FileNames []string
}

var searchPaths = []string{"./"}
var search = ".gcov" // look for these file extensions

// Search searches the designated paths for GCov files,
// appending them to the list of filenames.
func (f *Formatter) Search(paths ...string) (string, error) {
	paths = append(paths, searchPaths...)

	for _, p := range paths {
		logrus.Debugf("checking search path %s for GCov formatter", p)
		files, err := ioutil.ReadDir(p)
		if err != nil {
			return "", errors.WithStack(err)
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), search) {
				f.FileNames = append(f.FileNames, filepath.Join(p, file.Name()))
			}
		}
	}

	if len(f.FileNames) == 0 {
		return "", errors.WithStack(errors.Errorf("could not find any files in search paths for GCov. search paths were: %s", strings.Join(paths, ", ")))
	}

	return fmt.Sprint(f.FileNames), nil
}

// Format combines the source files into a report.
func (f *Formatter) Format() (formatters.Report, error) {
	rep, err := formatters.NewReport()
	if err != nil {
		return rep, err
	}

	for _, file := range f.FileNames {
		sf, err := parseSourceFile(file, nil)
		if err != nil {
			return rep, errors.WithStack(err)
		}
		err = rep.AddSourceFile(sf)
		if err != nil {
			return rep, errors.WithStack(err)
		}
	}

	return rep, nil
}

// Parse a single GCov source file.
func parseSourceFile(fileName string, gitHead *object.Commit) (formatters.SourceFile, error) {
	var sf formatters.SourceFile
	sourceFileName, err := getSourceFileName(fileName)
	if err != nil {
		return sf, errors.WithStack(err)
	}

	sf, err = formatters.NewSourceFile(sourceFileName, gitHead)
	if err != nil {
		return sf, errors.WithStack(err)
	}

	coverageFile, err := os.Open(fileName)
	if err != nil {
		return sf, errors.WithStack(err)
	}
	defer coverageFile.Close()

	scanner := bufio.NewScanner(coverageFile)

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.SplitN(string(line), ":", 3)
		if len(split) != 3 {
			continue
		}

		coverage := strings.TrimSpace(split[0])
		lineNum, _ := strconv.Atoi(strings.TrimSpace(split[1]))
		if lineNum < 1 { // pre code metadata
			continue
		}

		switch coverage {
		case "-":
			sf.Coverage = append(sf.Coverage, formatters.NullInt{})
		case "#####", "=====":
			sf.Coverage = append(sf.Coverage, formatters.NewNullInt(0))
		default: // coverage is number of hits
			// trailing * means that gcov detected an unexecuted block, and we don't
			// care that deeply. Example - a single-line if (x) { a } else { b } will
			// report with an asterisk if block b is never executed.
			if len(coverage) > 0 && coverage[len(coverage)-1] == '*' {
				coverage = coverage[:len(coverage)-1]
			}
			num, err := strconv.Atoi(coverage)
			if err != nil {
				return sf, errors.WithStack(err)
			}
			sf.Coverage = append(sf.Coverage, formatters.NewNullInt(num))
		}

	}

	return sf, nil
}

func getSourceFileName(coverageFileName string) (string, error) {
	coverageFile, err := os.Open(coverageFileName)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer coverageFile.Close()

	scanner := bufio.NewScanner(coverageFile)
	var sourceFileName string

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.SplitN(string(line), ":", 4)
		if len(split) != 4 {
			return "", errors.WithStack(errors.Errorf("Could not find source file name: %s", coverageFile.Name()))
		}
		sourceFileName = strings.TrimSpace(split[3])
		return sourceFileName, nil
	}

	return sourceFileName, err
}
