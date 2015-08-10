package testing

import (
	"github.com/synful/kudos/lib/yellow_dingo"
	"strings"
	"testing"
)

type assignmentErrTest struct {
	conf string
	err  string
}

var assignmentErrTests = []assignmentErrTest{
	{`{}`, "parse: bad code: must be nonempty"},
	{`{"code":""}`, "parse: bad code: must be nonempty"},
	{`{"code":"-"}`, "parse: bad code: contains illegal characters;" +
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
}

func TestParseAssignmentError(t *testing.T) {
	for _, test := range assignmentErrTests {
		_, err := yellow_dingo.ParseAssignment(strings.NewReader(test.conf))
		if err == nil || err.Error() != test.err {
			t.Errorf("unexpected error; want %v; got %v", test.err, err)
		}
	}
}
