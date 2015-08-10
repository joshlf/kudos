package purple_unicorn

import (
	"fmt"
	"time"
)

type Handin struct {
	code     *Code
	due      time.Time
	problems []Code
}

type Problem struct {
	code        Code
	name        string
	points      float64
	subproblems []Code
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
}

func (a *Assignment) Validate() error {
	if err := a.code.Validate(); err != nil {
		return fmt.Errorf("bad code: %v", err)
	}
	// validate problems first since handins refer to them,
	// and errors in the handins might actually derive from
	// errors in specifying problems

	switch len(a.handins) {
	case 0:
		return fmt.Errorf("must have at least one handin")
	case 1:
		if a.handins[0].code != nil {
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
			if h.code == nil {
				return fmt.Errorf("multiple handins defined; each must have a handin code")
			}
			if err := h.code.Validate(); err != nil {
				return fmt.Errorf("bad handin code %v: %v", *h.code, err)
			}
			hc := *h.code
			if seenCodes[*h.code] {
				return fmt.Errorf("duplicate handin code: %v", hc)
			}
			for _, pc := range h.problems {
				if err := pc.Validate(); err != nil {
					return fmt.Errorf("handin %v contains bad problem code %v: %v", hc, pc, err)
				}
				if _, ok := a.problemsByCode[pc]; !ok {
					return fmt.Errorf("handin %v specifies nonexistent problem: %v", hc, pc)
				}
				pu, ok := seenProblems[pc]
				if ok {
					if pu.problem == pc {
						return fmt.Errorf("handin %v includes problem %v, which was already included by handin %v", hc, pc, pu.handin)
					}
					return fmt.Errorf("handin %v includes problem %v, which was already included by handin %v as a subproblem of %v", hc, pc, pu.handin, pu.problem)
				}
				// rename so we aren't shadowed
				// by the argument to validate
				topLevelPC := pc
				var validate func(Code) error
				validate = func(pc Code) error {
					subproblems := a.problemsByCode[pc].subproblems
					for _, spc := range subproblems {
						pu, ok := seenProblems[spc]
						if ok {
							if pu.problem == spc {
								return fmt.Errorf("handin %v includes problem %v via %v, which was already included by handin %v", hc, spc, topLevelPC, pu.handin)
							}
							// If pu.problem != spc, that means that one of spc's
							// ancestors was included by another handin. However,
							// since we had to traverse down the tree to get here,
							// we should have already caught that.
							panic("internal error")
						}
						if err := validate(spc); err != nil {
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
			for _, pc = range a.problemsByCode[pc].subproblems {
				if err := validate(pc); err != nil {
					return err
				}
			}
			return nil
		}
		for _, p := range a.problems {
			if err := validate(p.code); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Assignment) MustValidate() {
	if err := a.Validate(); err != nil {
		panic(err)
	}
}
