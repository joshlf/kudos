package kudos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joshlf/kudos/lib/testutil"
)

var rubricErrTests = []struct {
	conf string
	err  string
}{
	{``, "EOF"},
	{`{}`, "must have uid or anonymous_token"},
	{`{"uid":"","anonymous_token":""}`, "cannot have both uid and anonymous_token"},
	{`{"uid":"a"}`, "must specify assignment"},
	{`{"uid":"","assignment":""}`, "uid must be numeric"},
	{`{"uid":"0","assignment":""}`, "bad assignment code \"\": must be non-empty"},
	{`{"anonymous_token":"","assignment":"a"}`, "anonymous_token must be hexadecimal"},
	{`{"anonymous_token":"a","assignment":"a"}`,
		"anonymous_token must have an even number of hexadecimal digits"},
	{`{"anonymous_token":"aa","assignment":"a","assignment":"a"}`,
		"must have at least one grade"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[]}`,
		"must have at least one grade"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{}]}`,
		"all grades must specify a problem"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":""}]}`,
		"bad problem code \"\": must be non-empty"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":"a"}]}`,
		"no grade given for problem a"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":"a","grade":""}]}`,
		"grade must be false or a number"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":"a","grade":true}]}`,
		"grade must be false or a number"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":"a","grade":false}]}`,
		"grade left unset for problem a"},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":"a","grade":0}]}`, ""},
	{`{"anonymous_token":"aa","assignment":"a","grades":[{"problem":"a","grade":0},
		{"problem":"a","grade":0}]}`, "duplicate grade for problem a"},
}

func TestParseRubricError(t *testing.T) {
	for i, test := range rubricErrTests {
		_, err := parseRubric(strings.NewReader(test.conf))
		prefix := fmt.Sprintf("test case %v (`%v`)", i, test.conf)
		if test.err == "" {
			testutil.MustPrefix(t, prefix, err)
		} else {
			testutil.MustErrorPrefix(t, prefix, test.err, err)
		}
	}
}
