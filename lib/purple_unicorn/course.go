package purple_unicorn

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Course struct {
	code         Code
	name         string
	description  string
	taGroup      *string
	tas          []string //TODO:  []TA?, how should this be represented
	studentGroup string
	handinMethod string
}

func (c *Course) Validate() error {
	if err := c.code.Validate(); err != nil {
		return err
	}
	if c.taGroup != nil && len(*c.taGroup) == 0 {
		return fmt.Errorf("Ta Group must not be the empty string; if you don't want a ta group, remove the line")
		//TODO: make the checking for group existing nicer, etc/group is not found on too many unices
		if f, err := os.Open("/etc/group"); err == nil {
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("Could not verify group: %v\n", buf)
			}
			str := string(buf)
			for _, line := range strings.Split(str, "\n") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 && parts[0] == *c.taGroup {
					goto success
				}
			}
			return fmt.Errorf("no such group: %v", *c.taGroup)

		success:
		}
	}
	return nil
}

//func (c *Course) SetTAGroup(group string)
