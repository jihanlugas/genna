package named

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

{{range $model := .Entities}}
type {{.GoName}} struct {
	{{range .Columns}}
	{{.GoName}} {{.Type}} {{.Tag}} {{.Comment}}{{end}}
}
{{end}}
`
