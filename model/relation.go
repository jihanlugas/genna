package model

import (
	"regexp"
	"strings"

	"github.com/dizzyfool/genna/util"
)

// Relation stores relation
type Relation struct {
	FKFields []string
	GoName   string

	TargetPGName     string
	TargetPGSchema   string
	TargetPGFullName string

	TargetEntity *Entity

	GoType string
}

// NewRelation creates relation from pg info
func NewRelation(sourceColumns []string, targetSchema, targetTable string) Relation {
	names := make([]string, len(sourceColumns))
	for i, name := range sourceColumns {
		names[i] = util.ReplaceSuffix(util.ColumnName(name), util.ID, "")
	}

	numRegEx := regexp.MustCompile(`[0-9]`)

	typ := util.EntityName(targetTable)
	typ = util.CamelCased(targetSchema) + typ
	typ = numRegEx.ReplaceAllString(typ, "")

	return Relation{
		FKFields: sourceColumns,
		GoName:   strings.Join(names, ""),

		TargetPGName:     targetTable,
		TargetPGSchema:   targetSchema,
		TargetPGFullName: util.JoinF(util.SchemaNameInFull(targetSchema), targetTable),
		// TargetPGFullName: util.JoinF(targetSchema, targetTable),

		GoType: typ,
	}
}

func (r *Relation) AddEntity(entity *Entity) {
	r.TargetEntity = entity
}
