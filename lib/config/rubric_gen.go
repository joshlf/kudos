package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Generate a default rubric given an assignment spec
func (a *AssignSpec) Rubric() (r Rubric) {
	res := Rubric{Assignment: a.Title, Grader: ""}

	for _, prob := range a.Problem {
		res.Grade = append(res.Grade, Grade{
			Problem:  prob.Name,
			Comment:  "",
			Score:    GradeNum{0},
			Possible: prob.Total,
		})
	}

	return res
}

// WARNING, MASSIVE JANK-FACTOR, WILL BE REMOVED ONCE TOML LIBRARY IS FIXED
func (r Rubric) WriteTOML(w io.Writer) error {
	var b bytes.Buffer

	enc := toml.NewEncoder(&b)

	if err := enc.Encode(r); err != nil {
		return err
	}

	for line, err := b.ReadString('\n'); ; line, err = b.ReadString('\n') {
		wout := line
		//fmt.Println("[tag]", wout)
		strs := strings.Split(line, "=")
		if len(strs) == 2 {
			if strings.Contains(strs[0], "score") || strings.Contains(strs[0], "possible") {
				wout = strings.Join([]string{strs[0],
					strings.Replace(strs[1], `"`, "", -1)}, "=")
			}

		}
		if _, err1 := fmt.Fprint(w, wout); err1 != nil {
			return err1
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (r Rubric) WriteTOMLFile(outFile string) error {
	f, err := os.Open(outFile)
	defer f.Close()

	if err != nil {
		return err
	}
	return r.WriteTOML(f)
}
