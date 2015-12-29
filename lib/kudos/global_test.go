package kudos

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshlf/kudos/lib/testutil"
)

var globalConfigErrTests = []struct {
	conf string
	err  string
}{
	{`{}`, "must have course_path_prefix"},
}

func TestParseGlobalConfigError(t *testing.T) {
	for _, test := range globalConfigErrTests {
		_, err := parseGlobalConfig(strings.NewReader(test.conf))
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("unexpected error; want %v; got %v", test.err, err)
		}
	}
}

func BenchmarkParseGlobalConfigFromDisk(b *testing.B) {
	dir, ok := testutil.SrcDir()
	if !ok {
		b.Skipf("could not determine source directory")
	}
	path := filepath.Join(dir, "testdata", "sample_global_config")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := ParseGlobalConfigFile(path)
		b.StopTimer()
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
