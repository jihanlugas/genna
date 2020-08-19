package model

import "github.com/dizzyfool/genna/model"

const (
	// Nil is nil check types
	Nil = "nil"
	// Zero is 0 check types
	Zero = "zero"
	// PZero is 0 check types for pointers
	PZero = "pzero"
	// Len is length check types
	Len = "len"
	// PLen is length check types for pointers
	PLen = "plen"
	// Enum is allowed values check types
	Enum = "enum"
	// PEnum is allowed values check types for pointers
	PEnum = "penum"
)

// isColumnValidatable checks if field can be validated
func isColumnValidatable(c model.Column) bool {
	// validate FK
	if c.IsFK {
		return true
	}

	// validate complex types
	if !c.Nullable && (c.IsArray || c.GoType == model.TypeMapInterface || c.GoType == model.TypeMapString) {
		return true
	}

	// validate strings len
	if c.GoType == model.TypeString && c.MaxLen > 0 {
		return true
	}

	// validate enum
	if len(c.Values) > 0 {
		return true
	}

	return false
}

// CheckColumn return string check type for validation
func CheckColumn(c model.Column) string {
	if !isColumnValidatable(c) {
		return ""
	}

	if c.IsArray || c.GoType == model.TypeMapInterface || c.GoType == model.TypeMapString {
		return Nil
	}

	if c.IsFK {
		if c.Nullable {
			return PZero
		}
		return Zero
	}

	if c.GoType == model.TypeString && c.MaxLen > 0 {
		if c.Nullable {
			return PLen
		}
		return Len
	}

	if len(c.Values) > 0 {
		if c.Nullable {
			return PEnum
		}
		return Enum
	}

	return ""
}
