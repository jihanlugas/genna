package model

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

{{range $model := .Entities}}
type {{.GoName}} struct {
	{{range .Columns}}
	{{.GoName}} {{.Type}} {{.Tag}} {{.Comment}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.GoName}} *{{.GoType}} {{.Tag}} {{.Comment}}{{end}}{{end}}
}
{{end}}
`
