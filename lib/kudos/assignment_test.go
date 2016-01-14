package kudos

import (
	"fmt"
	"path/filepath"
	"reflect"
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
		"handin b includes problem a, which is also included by handin a"},
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
	for i, test := range assignmentErrTests {
		_, err := parseAssignment(strings.NewReader(test.conf))
		prefix := fmt.Sprintf("test case %v (`%v`)", i, test.conf)
		if test.err == "" {
			testutil.MustPrefix(t, prefix, err)
		} else {
			testutil.MustErrorPrefix(t, prefix, test.err, err)
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

const findProblemPathByCodeTestAssignment = `{"code":"a","name":"a",
	"due":"Jul 4, 2015 at 12:00am (EST)","handins":[{"code":"first",
	"due":"Jul 4, 2015 at 12:00am (EST)","problems":["prob1"]},
	{"code":"second","due":"Jul 5, 2015 at 12:00am (EST)","problems":
	["prob2"]}],
	"problems": [
					{
						"code": "prob1",
						"name": "Problem 1",
						"points": 50
					},
					{
						"code": "prob2",
						"name": "Problem 2",
						"points": 50,
						"subproblems": [
										{
										"code": "a",
										"points": 25
										},
										{
										"code": "b",
										"points": 25
										}
										]
					}
				]
}`

var findProblemPathByCodeTestCases = []struct {
	code string
	path []string
	ok   bool
}{
	{"prob1", []string{}, true},
	{"prob2", []string{}, true},
	{"a", []string{"prob2"}, true},
	{"b", []string{"prob2"}, true},
	{"c", []string{}, false},
}

func TestFindProblemPathByCode(t *testing.T) {
	asgn, err := parseAssignment(strings.NewReader(findProblemPathByCodeTestAssignment))
	testutil.Must(t, err)
	for i, test := range findProblemPathByCodeTestCases {
		path, ok := asgn.FindProblemPathByCode(test.code)
		if len(path) == 0 {
			// make reflect.DeepEqual happy
			// (in case path is nil, which
			// reflect.DeepEqual will consider
			// unequal to an initialized,
			// zero-length slice)
			path = []string{}
		}
		prefix := fmt.Sprintf("test case %v (code %v)", i, test.code)
		if ok != test.ok || !reflect.DeepEqual(path, test.path) {
			t.Errorf("%v: got (%v, %v); want (%v, %v)", prefix, path, ok, test.path, test.ok)
		}
	}
}
