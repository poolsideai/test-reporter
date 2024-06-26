package formatters

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Report_Merge(t *testing.T) {
	r := require.New(t)
	reps := []*Report{}
	for i := 0; i < 4; i++ {
		rep, err := NewReport()
		r.NoError(err)

		f, err := os.Open(fmt.Sprintf("../integration-tests/codeclimate.%d.json", i))
		r.NoError(err)
		err = json.NewDecoder(f).Decode(&rep)
		r.NoError(err)

		sf := rep.SourceFiles["config/initializers/resque.rb"]
		r.NotNil(sf)
		r.Equal(14, sf.LineCounts.Total)

		reps = append(reps, &rep)
	}
	main := reps[0]
	main.Merge(reps[1:]...)
	r.Equal(19379, main.LineCounts.Total)
	r.Equal(2564, main.LineCounts.Missed)
	r.Equal(16815, main.LineCounts.Covered)
	r.InDelta(86.76, main.LineCounts.CoveredPercent(), 1)

	sf := main.SourceFiles["lib/code_climate/polymorphic_routes.rb"]
	r.NotNil(sf)
	r.Equal(196, sf.LineCounts.Total)
	r.Equal(59, sf.LineCounts.Missed)
	r.Equal(137, sf.LineCounts.Covered)
	r.InDelta(69.8, sf.CoveredPercent, 1)
}

func Test_Report_JSON_Unmarshal(t *testing.T) {
	r := require.New(t)
	f, err := os.Open("../integration-tests/codeclimate.json")
	r.NoError(err)

	rep, err := NewReport()
	r.NoError(err)
	err = json.NewDecoder(f).Decode(&rep)
	r.NoError(err)

	r.Equal(20, len(rep.SourceFiles))

	sf := rep.SourceFiles["lib/code_climate/test_reporter/client.rb"]
	r.NotNil(sf)
	r.InDelta(87.87, sf.CoveredPercent, 1)

	lc := sf.LineCounts
	r.Equal(8, lc.Missed)
	r.Equal(58, lc.Covered)
	r.Equal(66, lc.Total)
}

func Test_Merge_Issue_103(t *testing.T) {
	r := require.New(t)

	a, err := NewReport()
	r.NoError(err)

	sf := SourceFile{
		Name:           "app/jobs/initialize_account_seats.rb",
		CoveredPercent: 100,
		Coverage:       Coverage{NewNullInt(1), NewNullInt(1), NewNullInt(15), NullInt{}, NullInt{}, NewNullInt(1), NewNullInt(3), NullInt{}, NullInt{}, NewNullInt(1), NullInt{}, NewNullInt(1), NewNullInt(3), NullInt{}, NullInt{}},
		LineCounts: LineCounts{
			Missed:  0,
			Covered: 8,
			Total:   8,
		},
	}
	a.AddSourceFile(sf)

	b, err := NewReport()
	r.NoError(err)

	sf2 := SourceFile{
		Name:           "app/jobs/initialize_account_seats.rb",
		CoveredPercent: 62.5,
		Coverage:       Coverage{NewNullInt(1), NewNullInt(1), NewNullInt(0), NullInt{}, NullInt{}, NewNullInt(1), NewNullInt(0), NullInt{}, NullInt{}, NewNullInt(1), NullInt{}, NewNullInt(1), NewNullInt(0), NullInt{}, NullInt{}},
		LineCounts: LineCounts{
			Missed:  3,
			Covered: 5,
			Total:   8,
		},
	}

	b.AddSourceFile(sf2)

	err = a.Merge(&b)
	r.NoError(err)

	r.InDelta(100, a.CoveredPercent, 1)
	r.Equal(0, a.LineCounts.Missed)
	r.Equal(8, a.LineCounts.Covered)
	r.Equal(8, a.LineCounts.Total)
}
