package purple_unicorn

import (
	"strings"
	"testing"
)

type assignmentErrTest struct {
	conf string
	err  string
}

var assignmentErrTests = []assignmentErrTest{
	{`{}`, "parse: must have code"},
	{`{"code":""}`, "parse: bad code: must be nonempty"},
	{`{"code":"-"}`, "parse: bad code \"-\": contains illegal characters;" +
		" must be alphanumeric and start with an alphabetic character"},
	{`{"code":"a"}`, "parse: must have at least one problem"},
	{`{"code":"a","problems":[]}`, "parse: must have at least one problem"},
	{`{"code":"a","problems":[{}]}`, "parse: all problems must have codes"},
	{`{"code":"a","problems":[{"code":""}]}`,
		"parse: bad problem code : must be nonempty"},
	{`{"code":"a","problems":[{"code":"a"}]}`, "parse: problem a must have points"},
	{`{"code":"a","problems":[{"code":"a","points":1}]}`,
		"parse: must have at least one handin"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":[]}`,
		"parse: must have at least one handin"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":[{}]}`,
		"parse: handin must have due date"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)"}]}`,
		"parse: handin must specify at least one problem"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":[]}]}`,
		"parse: handin must specify at least one problem"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":[""]}]}`,
		"parse: handin contains bad problem code : must be nonempty"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["b"]}]}`,
		"parse: handin specifies nonexistent problem: b"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"parse: multiple handins defined; each must have a handin code"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"parse: multiple handins defined; each must have a handin code"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"parse: duplicate handin code: a"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]}]}`,
		"parse: handin b includes problem a, which was already included by handin a"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["b"]}]}`,
		"parse: handin b specifies nonexistent problem: b"},
	{`{"code":"a","problems":[{"code":"a","points":1}],"handins":
	[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":[]}]}`,
		"parse: handin b must specify at least one problem"},
	{`{"code":"a","problems":[{"code":"a","points":1},{"code":"a","points":1}]
	,"handins":[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":[]}]}`,
		"parse: duplicate problem code: a"},
	{`{"code":"a","problems":[{"code":"a","points":1,"subproblems":[{"code":"a"}]},
	{"code":"b","points":1}]
	,"handins":[{"code":"a","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["a"]},
	{"code":"b","due":"Jan 2, 2006 at 3:04pm (MST)","problems":["b"]}]}`,
		"parse: duplicate problem code: a"},
}

func TestParseAssignmentError(t *testing.T) {
	for _, test := range assignmentErrTests {
		_, err := ParseAssignment(strings.NewReader(test.conf))
		if err == nil || err.Error() != test.err {
			t.Errorf("unexpected error; want %v; got %v", test.err, err)
		}
	}
}
