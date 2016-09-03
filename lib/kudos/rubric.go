package kudos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type RubricGrade struct {
	Problem string
	Grade   float64
	Comment string
}

type Rubric struct {
	Anonymous           bool
	UID, AnonymousToken string

	Assignment string

	Grades []RubricGrade
}

func (r *Rubric) GetUID(ctx *Context) (uid string, ok bool) {
	if !r.Anonymous {
		return r.UID, true
	}
	return ctx.DB.Anonymizer.LookupToken(r.AnonymousToken)
}

// A jsonVerifiedGrade accepts a json field whose value
// is either a boolean or a number. If it is a boolean,
// then it is taken to mean that the original grade
// was not overwritten by the user.
type jsonVerifiedGrade struct {
	set   bool
	grade float64
}

func (j *jsonVerifiedGrade) UnmarshalJSON(b []byte) error {
	var empty bool
	err := json.Unmarshal(b, &empty)
	if err == nil {
		if empty {
			return fmt.Errorf("grade must be false or a number")
		}
		j.set = false
		return nil
	}

	var grade float64
	err = json.Unmarshal(b, &grade)
	if err != nil {
		return fmt.Errorf("grade must be false or a number")
	}
	j.set = true
	j.grade = grade
	return nil
}

func (j *jsonVerifiedGrade) MarshalJSON() ([]byte, error) {
	if j.set {
		return json.Marshal(j.grade)
	}
	return json.Marshal(false)
}

type parseableRubricGrade struct {
	Problem *string            `json:"problem"`
	Grade   *jsonVerifiedGrade `json:"grade"`
	Comment *string            `json:"comment"`
}

type parseableRubric struct {
	UID            *string `json:"uid,omitempty"`
	AnonymousToken *string `json:"anonymous_token,omitempty"`

	Assignment *string `json:"assignment"`

	Grades []parseableRubricGrade `json:"grades"`
}

func GenerateRubric(w io.Writer, asgn *Assignment, uid, token string, problems ...string) error {
	switch {
	case len(uid) != 0 && len(token) != 0:
		panic("kudos: both uid and token given")
	case len(uid) == 0 && len(token) == 0:
		panic("kudos: neither uid nor token given")
	}
	useToken := len(token) != 0

	var probs []Problem
	for _, p := range problems {
		pp, _ := asgn.FindProblemByCode(p)
		probs = append(probs, pp)
	}

	var r parseableRubric
	if useToken {
		r.AnonymousToken = &token
	} else {
		r.UID = &uid
	}
	// it's ok to use this address
	// because r won't be used after
	// the function returns
	r.Assignment = &asgn.Code

	for _, code := range problems {
		// make a local copy so that when
		// we take its address, it's unique
		// from the address of variables in
		// other loop iterations
		code := code
		p, _ := asgn.FindProblemByCode(code)
		r.Grades = append(r.Grades, parseableRubricGrade{
			Problem: &code,
			// the zero value is the empty string,
			// so it's fine to do this blindly
			Comment: &p.RubricCommentTemplate,
			Grade:   &jsonVerifiedGrade{},
		})
	}

	buf, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func ParseRubricFile(path string) (*Rubric, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r, err := parseRubric(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse: %v", err)
	}
	return r, nil
}

func parseRubric(r io.Reader) (*Rubric, error) {
	d := json.NewDecoder(r)
	var rubric parseableRubric
	err := d.Decode(&rubric)
	if err != nil {
		return nil, err
	}
	if err := validateRubric(rubric); err != nil {
		return nil, err
	}

	rr := &Rubric{
		Anonymous:  rubric.AnonymousToken != nil,
		Assignment: *rubric.Assignment,
	}
	if rr.Anonymous {
		rr.AnonymousToken = *rubric.AnonymousToken
	} else {
		rr.UID = *rubric.UID
	}

	for _, g := range rubric.Grades {
		gg := RubricGrade{
			Problem: *g.Problem,
			Grade:   g.Grade.grade,
		}
		if g.Comment != nil {
			gg.Comment = *g.Comment
		}
		rr.Grades = append(rr.Grades, gg)
	}
	return rr, nil
}

func validateRubric(r parseableRubric) error {
	switch {
	case r.AnonymousToken != nil && r.UID != nil:
		return fmt.Errorf("cannot have both uid and anonymous_token")
	case r.AnonymousToken == nil && r.UID == nil:
		return fmt.Errorf("must have uid or anonymous_token")
	case r.Assignment == nil:
		return fmt.Errorf("must specify assignment")
	}
	if r.UID != nil {
		numeric := true
		for _, c := range *r.UID {
			if !(c >= '0' && c <= '9') {
				numeric = false
			}
		}
		if !numeric || len(*r.UID) == 0 {
			return fmt.Errorf("uid must be numeric")
		}
	} else {
		if _, err := hex.DecodeString(*r.AnonymousToken); err != nil ||
			len(*r.AnonymousToken) == 0 {
			if err == hex.ErrLength {
				return fmt.Errorf("anonymous_token must have an even number of hexadecimal digits")
			}
			return fmt.Errorf("anonymous_token must be hexadecimal")
		}
	}

	if err := ValidateCode(*r.Assignment); err != nil {
		return fmt.Errorf("bad assignment code %q: %v", *r.Assignment, err)
	}

	if len(r.Grades) == 0 {
		return fmt.Errorf("must have at least one grade")
	}
	seenGrades := make(map[string]bool)
	for _, g := range r.Grades {
		if g.Problem == nil {
			return fmt.Errorf("all grades must specify a problem")
		}
		if err := ValidateCode(*g.Problem); err != nil {
			return fmt.Errorf("bad problem code %q: %v", *g.Problem, err)
		}
		if seenGrades[*g.Problem] {
			return fmt.Errorf("duplicate grade for problem %v", *g.Problem)
		}
		seenGrades[*g.Problem] = true
		switch {
		case g.Grade == nil:
			return fmt.Errorf("no grade given for problem %v", *g.Problem)
		case !g.Grade.set:
			return fmt.Errorf("grade left unset for problem %v", *g.Problem)
		}
	}
	return nil
}
