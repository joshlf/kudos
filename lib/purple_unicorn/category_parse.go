package purple_unicorn

import (
	"encoding/json"
	"fmt"
	"io"
)

type branchCategory struct {
	Name   *string `json:"category"`
	Weight *uint64 `json:"weight"`
}

func parseCategory(r io.Reader) (branchCategory, error) {
	var errs ErrList
	d := json.NewDecoder(r)
	var b branchCategory
	err := d.Decode(&b)
	if err != nil {
		return branchCategory{}, err
	}
	if b.Name == nil {
		errs.Add(fmt.Errorf("must add \"category\" field"))
	}
	if b.Weight == nil {
		errs.Add(fmt.Errorf("must add \"weight\" field"))
	}
	if len(errs) == 0 {
		return b, nil
	} else {
		return branchCategory{}, errs
	}
}
