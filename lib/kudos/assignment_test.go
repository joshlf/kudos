package kudos

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshlf/kudos/lib/testutil"
)

var assignmentErrTests = []struct {
	conf string
	err  string
}{
	{`{}`, "must have code"},
	{`{"code":""}`, "bad assignment code \"\": must be non-empty"},
	{`{"code":"-"}`, "bad assignment code \"-\": contains illegal characters;" +
		" must be alphanumeric and start with an alphabetic character"},
	{`{"code":"a"}`, "must have at least one problem"},
	{`{"code":"a","problems":[]}`, "must have at least one problem"},
	{`{"code":"a","problems":[{}]}`, "all problems must have codes"},
	{`{"code":"a","problems":[{"code":""}]}`,
		"bad problem code \"\": must be non-empty"},
	{`{"code":"a","problems":[{"code":"a"}]}`, "problem a must have points"},
	{`{"code":"a","problems":[{"code":"a","points":1}]}`,
		"must have at least one handin"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":[]}`,
		"must have at least one handin"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":[{}]}`,
		"handin must have due date"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)"}]}`,
		"handin must specify at least one problem"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":[]}]}`,
		"handin must specify at least one problem"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":[""]}]}`,
		"handin contains bad problem code \"\": must be non-empty"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["b"]}]}`,
		"handin specifies nonexistent problem: b"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"multiple handins defined; each must have a handin code"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"multiple handins defined; each must have a handin code"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"duplicate handin code: a"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"handin b includes problem a, which was already included by handin a"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["b"]}]}`,
		"handin b specifies nonexistent problem: b"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":[]}]}`,
		"handin b must specify at least one problem"},
	{`{"code":"a","problems":[{"code":"a","points":1},{"code":"a","points":1}]
	,"handins":[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":[]}]}`,
		"duplicate problem code: a"},
	{`{"code":"a","problems":[{"code":"a","points":1,"subproblems":[{"code":"a"}]},
	{"code":"b","points":1}]
	,"handins":[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["b"]}]}`,
		"duplicate problem code: a"},
	{`{"code":"a","problems":[{"code":"a","points":1,"subproblems":
	[{"code":"b","points":1},{"code":"c","points":1}]}],
	"handins":[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"problem a's points value is not equal to the sum of all subproblems' points"},
	{`{"code":"a","problems":[{"code":"a","points":2,"subproblems":
	[{"code":"b","points":1},{"code":"c","points":1}]}],
	"handins":[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		""},
}

func TestParseAssignmentError(t *testing.T) {
	for _, test := range assignmentErrTests {
		_, err := parseAssignment(strings.NewReader(test.conf))
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("unexpected error; want %v; got %v", test.err, err)
		}
	}
}

func BenchmarkParseAssignmentFromDisk(b *testing.B) {
	dir, ok := testutil.SrcDir()
	if !ok {
		b.Skipf("could not determine source directory")
	}
	path := filepath.Join(dir, "testdata", "sample_assignment")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		_, err := ParseAssignmentFile(path)
		b.StopTimer()
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
