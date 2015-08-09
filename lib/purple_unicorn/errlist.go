package purple_unicorn

import (
	"bytes"
	"text/template"
)

type ErrList []error

func (e ErrList) Error() string {
	b := bytes.NewBuffer([]byte{})
	template.Must(template.New("matcher").Parse(`
{{range .}}
    {{.Error()}}
{{end}}
`)).Execute(b, e)
	return b.String()
}
