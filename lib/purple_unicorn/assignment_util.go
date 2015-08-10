package purple_unicorn

import (
	"fmt"
	"time"
)

type problemInterface interface {
	code() Code
	name() string
	points() float64
	subproblems() []problemInterface
}

type parseableProblemInterface interface {
	problemInterface
	hasCode() bool
	hasName() bool
	hasPoints() bool
}

// validateProblemTree validates a tree of problems
// either for normal validation or parsing validation.
// This function assumes that the tree is homogeneous -
// either all of the problems only satisfy the problemInterface
// interface, or they all satisfy the parseableProblemInterface
// interface.
func validateProblemTree(problems []problemInterface) error {
	if len(problems) == 0 {
		return fmt.Errorf("must have at least one problem")
	}
	// From here on out, just use parseable
	// to determine the type, and assume that
	// type assertions will succeed
	_, parseable := problems[0].(parseableProblemInterface)

	// Check for code validity - valid codes and no duplicates
	seenCodes := make(map[Code]bool)
	var walkTree func(problems []problemInterface) error
	walkTree = func(problems []problemInterface) error {
		for _, p := range problems {
			if parseable && !p.(parseableProblemInterface).hasCode() {
				return fmt.Errorf("all problems must have codes")
			}
			c := p.code()
			switch {
			case c.Validate() != nil:
				return fmt.Errorf("bad problem code %v: %v", c, c.Validate())
			case seenCodes[c]:
				return fmt.Errorf("duplicate problem code: %v", c)
			}
			seenCodes[c] = true
			if err := walkTree(p.subproblems()); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walkTree(problems); err != nil {
		return err
	}

	// If these are parseable, also check that they all have points
	if parseable {
		walkTree = func(problems []problemInterface) error {
			for _, p := range problems {
				if !p.(parseableProblemInterface).hasPoints() {
					return fmt.Errorf("problem %v must have points", p.code())
				}
				if err := walkTree(p.subproblems()); err != nil {
					return err
				}
			}
			return nil
		}
		if err := walkTree(problems); err != nil {
			return err
		}
	}
	return nil
}

type handinInterface interface {
	code() Code
	hasCode() bool
	due() time.Time
	problems() []Code
}

type parseableHandinInterface interface {
	handinInterface
	hasDue() bool
}

// validateHandins validates a set of handins
// either for normal validation or parsing validation.
// This function assumes that the data is homogeneous -
// either all handins and problems satisfy the
// handinInterface and problemInterface interfaces, or
// they all satisfy the parseableHandinInterface and
// parseableProblemInterface interfaces.
func validateHandins(handins []handinInterface, problems []problemInterface) error {
	if len(handins) == 0 {
		return fmt.Errorf("must have at least one handin")
	}
	_, parseable := handins[0].(parseableHandinInterface)

	problemsByCode := make(map[Code]problemInterface)
	var walkProblems func(problems []problemInterface)
	walkProblems = func(problems []problemInterface) {
		for _, p := range problems {
			problemsByCode[p.code()] = p
			walkProblems(p.subproblems())
		}
	}
	walkProblems(problems)

	seenHandinCodes := make(map[Code]bool)
	// this type represents where a given
	// problem has been used before so we
	// can give helpful error messages when
	// problems are used by multiple handins
	type problemUsage struct {
		// the handin that used the problem
		handin Code

		// The problem that was directly included
		// by the handin. If problem is not equal
		// to the map key that corresponds to it,
		// it means that the map key is a child of
		// the problem that was included. Otherwise,
		// it is the problem itself. We must set
		// this variable in either case so we can
		// tell the difference and return precise
		// error messages.
		problem Code
	}
	seenProblems := make(map[Code]problemUsage)
	for _, h := range handins {
		switch {
		case len(handins) == 1 && h.hasCode():
			return fmt.Errorf("one handin defined; cannot have handin code")
		case len(handins) > 1 && !h.hasCode():
			return fmt.Errorf("multiple handins defined; each must have a handin code")
		case len(handins) > 1:
			if err := h.code().Validate(); err != nil {
				return fmt.Errorf("bad handin code %v: %v", h.code(), err)
			}
		case parseable && !h.(parseableHandinInterface).hasDue():
			if len(handins) == 1 {
				return fmt.Errorf("handin must have due date")
			}
			return fmt.Errorf("handin %v must have due date", h.code())
		}
		handinErrorName := "handin"
		if len(handins) > 1 {
			handinErrorName += " " + string(h.code())
		}
		if len(handins) > 1 && seenHandinCodes[h.code()] {
			return fmt.Errorf("duplicate handin code: %v", h.code())
		}
		seenHandinCodes[h.code()] = true
		if len(h.problems()) == 0 {
			return fmt.Errorf("%v must specify at least one problem", handinErrorName)
		}
		for _, pc := range h.problems() {
			if err := pc.Validate(); err != nil {
				return fmt.Errorf("%v contains bad problem code %v: %v", handinErrorName, pc, err)
			}
			if _, ok := problemsByCode[pc]; !ok {
				return fmt.Errorf("%v specifies nonexistent problem: %v", handinErrorName, pc)
			}
			pu, ok := seenProblems[pc]
			if ok {
				if pu.problem == pc {
					return fmt.Errorf("%v includes problem %v, which was already included by handin %v", handinErrorName, pc, pu.handin)
				}
				return fmt.Errorf("%v includes problem %v, which was already included by handin %v as a subproblem of %v", handinErrorName, pc, pu.handin, pu.problem)
			}
			// rename so we aren't shadowed
			// by the argument to validate
			topLevelPC := pc
			var validate func(Code) error
			validate = func(pc Code) error {
				seenProblems[pc] = problemUsage{h.code(), topLevelPC}
				subproblems := problemsByCode[pc].subproblems()
				for _, sp := range subproblems {
					pu, ok := seenProblems[sp.code()]
					if ok {
						if pu.problem == sp.code() {
							return fmt.Errorf("%v includes problem %v via %v, which was already included by handin %v", handinErrorName, sp.code(), topLevelPC, pu.handin)
						}
						// If pu.problem != spc, that means that one of spc's
						// ancestors was included by another handin. However,
						// since we had to traverse down the tree to get here,
						// we should have already caught that.
						panic("internal error")
					}
					if err := validate(sp.code()); err != nil {
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
	// problemsByCode) so that we encounter errors
	// in order
	var validate func(Code) error
	validate = func(pc Code) error {
		if _, ok := seenProblems[pc]; !ok {
			return fmt.Errorf("problem %v not in any handins", pc)
		}
		for _, sp := range problemsByCode[pc].subproblems() {
			if err := validate(sp.code()); err != nil {
				return err
			}
		}
		return nil
	}
	for _, p := range problems {
		if err := validate(p.code()); err != nil {
			return err
		}
	}
	return nil
}

// validateProblemsAndHandins validates a set of
// handins and a problem tree either for normal
// validation or parsing validation. This function
// assumes that the data is homogeneous - either all
// handins and problems satisfy the handinInterface
// and problemInterface interfaces, or they all
// satisfy the parseableHandinInterface and
// parseableProblemInterface interfaces.
func validateProblemsAndHandins(handins []handinInterface, problems []problemInterface) error {
	if err := validateProblemTree(problems); err != nil {
		return err
	}
	return validateHandins(handins, problems)
}
