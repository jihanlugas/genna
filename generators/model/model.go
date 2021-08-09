package model

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// TemplatePackage stores package info
type TemplatePackage struct {
	Package string

	HasImports bool
	Imports    []string
	HasEnums   bool
	Enums      []util.Enum

	Entities []TemplateEntity
}

// NewTemplatePackage creates a package for template
func NewTemplatePackage(entities []model.Entity, options Options) TemplatePackage {
	imports := util.NewSet()
	enums := util.NewSetEnum()

	models := make([]TemplateEntity, len(entities))
	for i, entity := range entities {
		for _, imp := range entity.Imports {
			imports.Add(imp)
		}

		for _, enm := range entity.Enums {
			enums.Add(enm)
		}

		models[i] = NewTemplateEntity(entity, options)
	}

	// imports.Add("time")

	return TemplatePackage{
		Package: options.Package,

		HasImports: imports.Len() > 0,
		Imports:    imports.Elements(),
		HasEnums:   enums.Len() > 0,
		Enums:      enums.Elements(),

		Entities: models,
	}
}

// TemplateEntity stores struct info
type TemplateEntity struct {
	model.Entity

	Tag template.HTML

	NoAlias bool
	Alias   string

	Columns []TemplateColumn

	HasRelations bool
	Relations    []TemplateRelation

	HasCreateBy  bool
	HasCreateDt  bool
	HasUpdateBy  bool
	HasUpdateDt  bool
	HasArchiveBy bool
	HasArchiveDt bool
}

// NewTemplateEntity creates an entity for template
func NewTemplateEntity(entity model.Entity, options Options) TemplateEntity {
	if entity.HasMultiplePKs() {
		options.KeepPK = true
	}

	columns := make([]TemplateColumn, len(entity.Columns))
	hasCreateBy := false
	hasCreateDt := false
	hasUpdateDt := false
	hasUpdateBy := false
	hasArchiveBy := false
	hasArchiveDt := false

	for i, column := range entity.Columns {
		switch column.GoName {
		case "CreateBy":
			hasCreateBy = true
		case "CreateDt":
			hasCreateDt = true
		case "UpdateDt":
			hasUpdateDt = true
		case "UpdateBy":
			hasUpdateBy = true
		case "ArchiveBy":
			hasArchiveBy = true
		case "ArchiveDt":
			hasArchiveDt = true
		}
		columns[i] = NewTemplateColumn(entity, column, options)
	}

	relations := make([]TemplateRelation, len(entity.Relations))
	for i, relation := range entity.Relations {
		relations[i] = NewTemplateRelation(relation, options)
	}

	tagName := tagName(options)
	tags := util.NewAnnotation()

	// fmt.Println(entity.GoName)
	// fmt.Println(entity.GoNamePlural)
	// fmt.Println(entity.PGName)
	// fmt.Println(entity.PGSchema)
	// fmt.Println(entity.PGFullName)

	// tags.AddTag(tagName, util.Quoted(entity.PGSchema+"."+entity.PGFullName, true))
	tags.AddTag(tagName, util.Quoted(entity.PGFullName, true))
	if !options.NoAlias {
		tags.AddTag(tagName, fmt.Sprintf("alias:%s", entity.PGName))
	}

	if !options.NoDiscard {
		// leading comma is required
		tags.AddTag("db", ",discard_unknown_columns")
	}

	return TemplateEntity{
		Entity: entity,
		Tag:    template.HTML(fmt.Sprintf("`%s`", tags.String())),

		NoAlias: options.NoAlias,
		Alias:   entity.PGName,

		Columns: columns,

		HasRelations: len(relations) > 0,
		Relations:    relations,

		HasCreateBy:  hasCreateBy,
		HasCreateDt:  hasCreateDt,
		HasUpdateBy:  hasUpdateBy,
		HasUpdateDt:  hasUpdateDt,
		HasArchiveBy: hasArchiveBy,
		HasArchiveDt: hasArchiveDt,
	}
}

// TemplateColumn stores column info
type TemplateColumn struct {
	model.Column

	Tag     template.HTML
	Comment template.HTML
}

// NewTemplateColumn creates a column for template
func NewTemplateColumn(entity model.Entity, column model.Column, options Options) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	if column.PGType == model.TypePGJSON || column.PGType == model.TypePGJSONB {
		if typ, ok := jsonType(options.JSONTypes, entity.PGSchema, entity.PGName, column.PGName); ok {
			column.Type = typ
		}
	}

	comment := ""
	tagName := tagName(options)
	tags := util.NewAnnotation()
	tags.AddTag(tagName, column.PGName)

	// pk tag
	if column.IsPK {
		tags.AddTag(tagName, "pk")
	}

	// types tag
	if column.PGType == model.TypePGHstore {
		tags.AddTag(tagName, "hstore")
	} else if column.IsArray {
		tags.AddTag(tagName, "array")
	}
	if column.PGType == model.TypePGUuid {
		tags.AddTag(tagName, "type:uuid")
	}

	// nullable tag
	if !column.Nullable && !column.IsPK {
		if options.GoPgVer == 9 {
			tags.AddTag(tagName, "use_zero")
		} else {
			tags.AddTag(tagName, "notnull")
		}
	}

	// soft_delete tag
	if options.SoftDelete == column.PGName && column.Nullable && column.GoType == model.TypeTime && !column.IsArray {
		tags.AddTag("db", ",soft_delete")
	}

	// ignore tag
	if column.GoType == model.TypeInterface {
		comment = "// unsupported"
		tags = util.NewAnnotation().AddTag(tagName, "-")
	}

	tags.AddTag("json", util.LowerFirst(util.ReplaceSuffix(util.ReplaceSuffix(column.GoName, util.ID, util.Id), util.IDs, util.Ids)))
	tags.AddTag("form", util.LowerFirst(util.ReplaceSuffix(util.ReplaceSuffix(column.GoName, util.ID, util.Id), util.IDs, util.Ids)))

	// if column.GoType == model.TypeInt64 {
	// 	tags.AddTag("json", "string")
	// }

	if !column.Nullable {
		tags.AddTag("validate", "required")
	}

	// validate complex types
	// if !column.Nullable && (column.IsArray || column.GoType == model.TypeMapInterface || column.GoType == model.TypeMapString) {
	// untuk validate sqltype hstore atau json
	// }

	// validate FK
	// if column.IsFK {
	// 	if column.Nullable {
	// 		return PZero
	// 	}
	// 	return Zero
	// }

	// validate enum
	if len(column.Values) > 0 {
		if column.Nullable {
			tags.AddTag("validate", "omitempty")
		}
		if column.IsArray {
			tags.AddTag("validate", "dive")
		}
		tags.AddTag("validate", "oneof="+fmt.Sprintf(`'%s'`, strings.Join(column.Values, `' '`)))
	}

	// validate strings len
	if column.GoType == model.TypeString {
		if column.MaxLen > 0 {
			tags.AddTag("validate", "lte="+strconv.Itoa(column.MaxLen))
		}
	}

	return TemplateColumn{
		Column: column,

		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
		Comment: template.HTML(comment),
	}
}

// TemplateRelation stores relation info
type TemplateRelation struct {
	model.Relation

	Tag     template.HTML
	Comment template.HTML
}

// NewTemplateRelation creates relation for template
func NewTemplateRelation(relation model.Relation, options Options) TemplateRelation {
	comment := ""
	tagName := tagName(options)
	tags := util.NewAnnotation().AddTag("db", "fk:"+strings.Join(relation.FKFields, ","))
	if len(relation.FKFields) > 1 {
		comment = "// unsupported"
		tags.AddTag(tagName, "-")
	}

	tags.AddTag("json", "-")

	return TemplateRelation{
		Relation: relation,

		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
		Comment: template.HTML(comment),
	}
}

func jsonType(mp map[string]string, schema, table, field string) (string, bool) {
	if mp == nil {
		return "", false
	}

	patterns := [][3]string{
		{schema, table, field},
		{schema, "*", field},
		{schema, table, "*"},
		{schema, "*", "*"},
	}

	var names []string
	for _, parts := range patterns {
		names = append(names, fmt.Sprintf("%s.%s", util.Join(parts[0], parts[1]), parts[2]))
		names = append(names, fmt.Sprintf("%s.%s", util.JoinF(parts[0], parts[1]), parts[2]))
	}
	names = append(names, util.Join(schema, table), "*")

	for _, name := range names {
		if v, ok := mp[name]; ok {
			return v, true
		}
	}

	return "", false
}

func tagName(options Options) string {
	return "db"
}
