package named

const EnumTemplate = `//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package constant{{if .HasEnums}}

const ({{range .Enums}}{{range .Entries}}
	{{.TagName}} = "{{.Value}}"{{end}}{{end}}
){{end}}
`
