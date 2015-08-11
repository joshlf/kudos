package purple_unicorn

import (
	"bytes"
	"text/template"
)

type ErrList []error

func (e ErrList) Error() string {
	var b bytes.Buffer
	template.Must(template.New("matcher").Parse(`{{range .}}{{.Error}}
{{end}}`)).Execute(&b, e)
	return b.String()
}

func (e *ErrList) Add(err error) {
	*e = append(*e, err)
}
