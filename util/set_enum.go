package util

type EnumEntries struct {
	TagName string
	Value   string
}

type Enum struct {
	Name    string
	Values  []string
	Entries []EnumEntries
}

// SetEnum stores only unique model.Enum
type SetEnum struct {
	elements []Enum
	index    map[string]struct{}
}

// NewSetEnum creates SetEnum
func NewSetEnum() SetEnum {
	return SetEnum{
		elements: []Enum{},
		index:    map[string]struct{}{},
	}
}

// Add adds element to set
// return false if element already exists
func (s *SetEnum) Add(element Enum) bool {
	if s.Exists(element.Name) {
		return false
	}

	vals := make([]EnumEntries, len(element.Values))
	for i, val := range element.Values {
		vals[i] = EnumEntries{
			TagName: CamelCased(element.Name + "_" + PackageName(val)),
			Value:   val,
		}
	}

	element.Entries = vals

	s.elements = append(s.elements, element)
	s.index[element.Name] = struct{}{}

	return true
}

// Exists checks if element exists
func (s *SetEnum) Exists(element string) bool {
	_, ok := s.index[element]
	return ok
}

// Elements return all elements from set
func (s *SetEnum) Elements() []Enum {
	return s.elements
}

// Len gets elements count
func (s *SetEnum) Len() int {
	return len(s.elements)
}
