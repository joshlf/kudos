package purple_unicorn

import (
	"fmt"
	"time"
)

// We never want users of this package to be able
// to interact with internal Assignment state other
// than through method calls. Thus, when consuming
// or producing Problems or Handins, make sure to use
// copyProblemTree or copyHandin to deep copy the
// Problem and all of its subproblems or the Handin's
// Problems slice. This avoids referencing state that
// is accessible from outside the package.

type Handin struct {
	Code     Code
	HasCode  bool
	Due      time.Time
	Problems []Code
}

// Handin implements the handinInterface interface
func (h Handin) code() Code     { return h.Code }
func (h Handin) hasCode() bool  { return h.HasCode }
func (h Handin) due() time.Time { return h.Due }

// At the time of writing, no caller of the problems method
// ever modifies the return value, so this deep copy isn't
// strictly necessary. However, having it anyway means that
// we avoid a potentially very subtle bug in the future.
func (h Handin) problems() []Code { return append([]Code(nil), h.Problems...) }

type Problem struct {
	Code        Code
	Name        string
	Points      float64
	Subproblems []Problem
}

// EffectiveName returns the effective name of p,
// which is p.Name if it is not empty, otherwise
// p.Code. EffectiveName should always be used
// unless the distinction between effective name
// and the name set by the user is important.
func (p Problem) EffectiveName() string {
	if p.Name != "" {
		return p.Name
	}
	return string(p.Code)
}

// Problem implements the problemInterface interface
func (p Problem) code() Code      { return p.Code }
func (p Problem) name() string    { return p.Name }
func (p Problem) points() float64 { return p.Points }
func (p Problem) subproblems() []problemInterface {
	var sp []problemInterface
	for _, pp := range p.Subproblems {
		sp = append(sp, pp)
	}
	return sp
}

type Assignment struct {
	code *Code
	name string

	// Only handins with codes
	// will be in handinsByCode
	handins       []Handin
	handinsByCode map[Code]Handin

	// Only contains top-level problems
	problems []Problem
	// Contains all problems and subproblems
	problemsByCode map[Code]Problem
}

func NewAssignment() *Assignment {
	return &Assignment{
		handinsByCode:  make(map[Code]Handin),
		problemsByCode: make(map[Code]Problem),
	}
}

func (a *Assignment) SetCodeNoValidate(c Code) {
	a.code = &c
}

func (a *Assignment) SetCode(c Code) error {
	if err := c.Validate(); err != nil {
		return err
	}
	a.code = &c
	return nil
}

func (a *Assignment) MustSetCode(c Code) {
	if err := a.SetCode(c); err != nil {
		mustPanic(err, "Assignment.SetCode")
	}
}

func (a *Assignment) Code() Code { return *a.code }

func (a *Assignment) SetName(name string) { a.name = name }
func (a *Assignment) Name() string        { return a.name }

// EffectiveName returns the effective name of a,
// which is a.Name() if it is not empty, otherwise
// a.Code(). EffectiveName should always be used
// unless the distinction between effective name
// and the name set by the user is important.
func (a *Assignment) EffectiveName() string {
	if name := a.Name(); name != "" {
		return name
	}
	return string(a.Code())
}

func (a *Assignment) SetProblemsNoValidate(p []Problem) {
	a.problems = copyProblemsTree(p)
	a.problemsByCode = make(map[Code]Problem)
	f := func(p Problem) {
		a.problemsByCode[p.Code] = p
	}
	mapProblemsTree(a.problems, f)
}

func (a *Assignment) SetProblems(p []Problem) error {
	problems := a.problems
	problemsByCode := a.problemsByCode
	a.SetProblemsNoValidate(p)
	if err := a.Validate(); err != nil {
		a.problems = problems
		a.problemsByCode = problemsByCode
		return err
	}
	return nil
}

func (a *Assignment) MustSetProblems(p []Problem) {
	if err := a.SetProblems(p); err != nil {
		mustPanic(err, "Assignment.SetProblems")
	}
}

func (a *Assignment) Problems() []Problem { return copyProblemsTree(a.problems) }
func (a *Assignment) ProblemByCode(code Code) (p Problem, ok bool) {
	p, ok = a.problemsByCode[code]
	if !ok {
		return
	}
	p.Subproblems = copyProblemsTree(p.Subproblems)
	return
}

// deep-copies p
func copyProblemsTree(p []Problem) []Problem {
	p = append([]Problem(nil), p...)
	for i := range p {
		p[i].Subproblems = copyProblemsTree(p[i].Subproblems)
	}
	return p
}

// maps f over p in a pre-order traversal
func mapProblemsTree(p []Problem, f func(p Problem)) {
	for _, pp := range p {
		f(pp)
		mapProblemsTree(pp.Subproblems, f)
	}
}

func (a *Assignment) SetHandinsNoValidate(h []Handin) {
	for i := range h {
		h[i] = copyHandin(h[i])
	}
	a.handins = h
	a.handinsByCode = make(map[Code]Handin)
	for _, hh := range h {
		if hh.HasCode {
			a.handinsByCode[hh.Code] = hh
		}
	}
}

func (a *Assignment) SetHandins(h []Handin) error {
	handins := a.handins
	handinsByCode := a.handinsByCode
	a.SetHandinsNoValidate(h)
	if err := a.Validate(); err != nil {
		a.handins = handins
		a.handinsByCode = handinsByCode
		return err
	}
	return nil
}

func (a *Assignment) MustSetHandins(h []Handin) {
	if err := a.SetHandins(h); err != nil {
		mustPanic(err, "Assignment.SetHandins")
	}
}

func (a *Assignment) Handins() []Handin { return copyHandins(a.handins) }

func (a *Assignment) HandinByCode(c Code) (h Handin, ok bool) {
	h, ok = a.handinsByCode[c]
	if !ok {
		return
	}
	h = copyHandin(h)
	return
}

// deep-copies h
func copyHandin(h Handin) Handin {
	h.Problems = append([]Code(nil), h.Problems...)
	return h
}

func copyHandins(h []Handin) []Handin {
	h = append([]Handin(nil), h...)
	for i := range h {
		h[i] = copyHandin(h[i])
	}
	return h
}

// Copy deep copies a, returning an identical
// but distinct Assignment which shares no
// underlying resources with the original;
// modifications to either will not affect the
// other.
func (a *Assignment) Copy() *Assignment {
	code := *a.code
	aa := &Assignment{
		code:           &code,
		name:           a.name,
		handins:        copyHandins(a.handins),
		handinsByCode:  make(map[Code]Handin),
		problems:       copyProblemsTree(a.problems),
		problemsByCode: make(map[Code]Problem),
	}
	for k, v := range a.handinsByCode {
		aa.handinsByCode[k] = v
	}
	for k, v := range a.problemsByCode {
		v.Subproblems = copyProblemsTree(v.Subproblems)
		aa.problemsByCode[k] = v
	}
	return aa
}

// Validate implements the Validator Validate method.
func (a *Assignment) Validate() error {
	if a.code == nil {
		return fmt.Errorf("must have code")
	}
	if err := a.code.Validate(); err != nil {
		return codeErrMsg(*a.code, err)
	}
	// validate problems first since validateHandins
	// requires that problems have already been validated
	var problems []problemInterface
	for _, p := range a.problems {
		problems = append(problems, p)
	}
	if err := validateProblemTree(problems); err != nil {
		return err
	}
	var handins []handinInterface
	for _, h := range a.handins {
		handins = append(handins, h)
	}
	return validateHandins(handins, problems)
}

// MustValidate implements the Validator MustValidate method.
func (a *Assignment) MustValidate() {
	if err := a.Validate(); err != nil {
		mustPanic(err, "Assignment.Validate")
	}
}
