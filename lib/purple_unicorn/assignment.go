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

type Assignment struct {
	code Code
	name string

	// Only handins with codes
	// will be in handinsByCode
	handins       []Handin
	handinsByCode map[Code]Handin

	// Only contains top-level problems
	problems []Problem
	// Contains all problems and subproblems
	problemsByCode map[Code]Problem

	// Able to store duplicate codes
	// (since problemsByCode can't)
	// so that we can catch errors
	// during validation
	allProblemCodes []Code
}

func (a *Assignment) SetCodeNoValidate(c Code) {
	a.code = c
}

func (a *Assignment) SetCode(c Code) error {
	if err := c.Validate(); err != nil {
		return err
	}
	a.code = c
	return nil
}

func (a *Assignment) MustSetCode(c Code) {
	if err := a.SetCode(c); err != nil {
		mustPanic(err, "Assignment.SetCode")
	}
}

func (a *Assignment) Code() Code { return a.code }

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
	a.allProblemCodes = nil
	f := func(p Problem) {
		a.problemsByCode[p.Code] = p
		a.allProblemCodes = append(a.allProblemCodes, p.Code)
	}
	mapProblemsTree(a.problems, f)
}

func (a *Assignment) SetProblems(p []Problem) error {
	problems := a.problems
	problemsByCode := a.problemsByCode
	allProblemCodes := a.allProblemCodes
	a.SetProblemsNoValidate(p)
	if err := a.Validate(); err != nil {
		a.problems = problems
		a.problemsByCode = problemsByCode
		a.allProblemCodes = allProblemCodes
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

func (a *Assignment) Handins() []Handin {
	h := append([]Handin(nil), a.handins...)
	for i := range h {
		h[i] = copyHandin(h[i])
	}
	return h
}

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

// Validate implements the Validator Validate method.
func (a *Assignment) Validate() error {
	if err := a.code.Validate(); err != nil {
		return fmt.Errorf("bad code: %v", err)
	}
	// validate problems first since handins refer to them,
	// and errors in the handins might actually derive from
	// errors in specifying problems
	seenCodes := make(map[Code]bool)
	for _, pc := range a.allProblemCodes {
		if err := pc.Validate(); err != nil {
			return fmt.Errorf("bad problem code %v: %v", pc, err)
		}
		if seenCodes[pc] {
			return fmt.Errorf("duplicate problem code: %v")
		}
	}
	// TODO(synful): validate points

	switch len(a.handins) {
	case 0:
		return fmt.Errorf("must have at least one handin")
	case 1:
		if a.handins[0].HasCode {
			return fmt.Errorf("one handin defined; cannot have handin code")
		}
	default:
		seenCodes := make(map[Code]bool)
		type problemUsage struct {
			// the handin that used the problem
			handin Code

			// the problem that was directly included;
			// if this is a subproblem, it will be a
			// parent or grandparent or higher; if this
			// is the problem itself, this will be
			// redundant (but must be set so we can check
			// against it)
			problem Code
		}
		seenProblems := make(map[Code]problemUsage)
		for _, h := range a.handins {
			if !h.HasCode {
				return fmt.Errorf("multiple handins defined; each must have a handin code")
			}
			if err := h.Code.Validate(); err != nil {
				return fmt.Errorf("bad handin code %v: %v", h.Code, err)
			}
			if seenCodes[h.Code] {
				return fmt.Errorf("duplicate handin code: %v", h.Code)
			}
			for _, pc := range h.Problems {
				if err := pc.Validate(); err != nil {
					return fmt.Errorf("handin %v contains bad problem code %v: %v", h.Code, pc, err)
				}
				if _, ok := a.problemsByCode[pc]; !ok {
					return fmt.Errorf("handin %v specifies nonexistent problem: %v", h.Code, pc)
				}
				pu, ok := seenProblems[pc]
				if ok {
					if pu.problem == pc {
						return fmt.Errorf("handin %v includes problem %v, which was already included by handin %v", h.Code, pc, pu.handin)
					}
					return fmt.Errorf("handin %v includes problem %v, which was already included by handin %v as a subproblem of %v", h.Code, pc, pu.handin, pu.problem)
				}
				// rename so we aren't shadowed
				// by the argument to validate
				topLevelPC := pc
				var validate func(Code) error
				validate = func(pc Code) error {
					subproblems := a.problemsByCode[pc].Subproblems
					for _, sp := range subproblems {
						pu, ok := seenProblems[sp.Code]
						if ok {
							if pu.problem == sp.Code {
								return fmt.Errorf("handin %v includes problem %v via %v, which was already included by handin %v", h.Code, sp.Code, topLevelPC, pu.handin)
							}
							// If pu.problem != spc, that means that one of spc's
							// ancestors was included by another handin. However,
							// since we had to traverse down the tree to get here,
							// we should have already caught that.
							panic("internal error")
						}
						if err := validate(sp.Code); err != nil {
							return err
						}
					}
					return nil
				}
				err := validate(topLevelPC)
				if err != nil {
					return err
				}
			}
		}
		// Traverse the problem tree (as opposed to using
		// a.problemsByCode) so that we encounter errors
		// in order
		var validate func(Code) error
		validate = func(pc Code) error {
			if _, ok := seenProblems[pc]; !ok {
				return fmt.Errorf("problem %v not in any handins", pc)
			}
			for _, sp := range a.problemsByCode[pc].Subproblems {
				if err := validate(sp.Code); err != nil {
					return err
				}
			}
			return nil
		}
		for _, p := range a.problems {
			if err := validate(p.Code); err != nil {
				return err
			}
		}
	}
	return nil
}

// MustValidate implements the Validator MustValidate method.
func (a *Assignment) MustValidate() {
	if err := a.Validate(); err != nil {
		mustPanic(err, "Assignment.Validate")
	}
}
