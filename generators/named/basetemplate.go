package named

const BaseTemplate = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}

type ColumnsSt struct { {{range .Entities}}
	{{.GoName}} Columns{{.GoName}}{{end}}
}
var Columns = ColumnsSt{ {{range .Entities}}
	{{.GoName}}: Columns{{.GoName}}{ {{range .Columns}}
		{{.GoName}}: "{{.PGName}}",{{end}}{{if .HasRelations}}
		{{range .Relations}}
		{{.GoName}}: "{{.GoName}}",{{end}}{{end}}
	},{{end}}
}
type TablesSt struct { {{range .Entities}}
		{{.GoName}} Table{{.GoName}}{{end}}
}
var Tables = TablesSt { {{range .Entities}}
	{{.GoName}}: Table{{.GoName}}{ 
		Name: "{{.PGFullName}}"{{if not .NoAlias}},
		Alias: "{{.Alias}}",{{end}}
	},{{end}}
}
`
