package kudos

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshlf/kudos/lib/testutil"
)

var courseErrTests = []struct {
	conf string
	err  string
}{
	{`{}`, "must have code"},
	{`{"code":""}`, "bad course code \"\": must be non-empty"},
	{`{"code":"-"}`, "bad course code \"-\": contains illegal characters;" +
		" must be alphanumeric and start with an alphabetic character"},
	{`{"code":"course"}`, "must have TA group"},
}

func TestParseCourseError(t *testing.T) {
	for _, test := range courseErrTests {
		_, err := parseCourse(strings.NewReader(test.conf))
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("unexpected error; want %v; got %v", test.err, err)
		}
	}
}

func BenchmarkParseCourseFromDisk(b *testing.B) {
	dir, ok := testutil.SrcDir()
	if !ok {
		b.Skipf("could not determine source directory")
	}
	path := filepath.Join(dir, "testdata", "sample_course")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := ParseCourseFile(path)
		b.StopTimer()
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
