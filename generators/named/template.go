package named

const Template = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

{{range .Entities}}
	type Columns{{.GoName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{if .HasRelations}}
		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.GoName}}{{end}} string{{end}}
	}
{{end}}

{{range .Entities}}
type Table{{.GoName}} struct {
	Name{{if not .NoAlias}}, Alias{{end}} string
}
{{end}}

{{range $model := .Entities}}
type {{.GoName}} struct {
	tableName struct{} {{.Tag}}
	{{range .Columns}}
	{{.GoName}} {{.Type}} {{.Tag}} {{.Comment}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.GoName}} *{{.GoType}} {{.Tag}} {{.Comment}}{{end}}{{end}}
}

func (m *{{.GoName}}) Name() string {
	return "{{.GoName}}"
}

func (m *{{.GoName}}) BeforeInsert(u IntStr, now *time.Time) {
	{{if .HasCreateBy}}m.CreateBy = u{{end}}{{if .HasCreateDt}}
	m.CreateDt = now{{end}}{{if .HasUpdateBy}}
	m.UpdateBy = u{{end}}{{if .HasUpdateDt}}
	m.UpdateDt = now{{end}}
}

func (m *{{.GoName}}) BeforeUpdate(u IntStr, now *time.Time) {
	{{if .HasUpdateBy}}m.UpdateBy = u{{end}}{{if .HasUpdateDt}}
	m.UpdateDt = now{{end}}
}

func (m *{{.GoName}}) BeforeArchive(u IntStr, now *time.Time) {
	{{if .HasArchiveBy}}m.ArchiveBy = u{{end}}{{if .HasArchiveDt}}
	m.ArchiveDt = now{{end}}{{if .HasUpdateBy}}
	m.UpdateBy = u{{end}}{{if .HasUpdateDt}}
	m.UpdateDt = now{{end}}
}
{{end}}
`
