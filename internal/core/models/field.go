package models

type FieldType string

const (
	DateType    FieldType = "date"
	IPType      FieldType = "ip"
	KeywordType FieldType = "keyword"
	TextType    FieldType = "text"
	IntegerType FieldType = "integer"
	LongType    FieldType = "long"
)

type Field struct {
	Type         FieldType `json:"type"`
	Searchable   bool      `json:"searchable"`
	Aggregatable bool      `json:"aggregatable"`
	Indices      []string  `json:"indices"`
}

// NewField creates a new Field struct
func NewField(fieldType FieldType, searchable, aggregatable bool) *Field {
	return &Field{
		Type:         fieldType,
		Searchable:   searchable,
		Aggregatable: aggregatable,
		Indices:      []string{},
	}
}

// Fields is a struct that represents the fields of an index
type Fields struct {
	Indices []string                        `json:"indices"`
	Fields  map[string]map[FieldType]*Field `json:"fields"`
}

// AddField adds a field to the Fields struct
func (f *Fields) AddField(name string, field *Field) {
	if f.Fields == nil {
		f.Fields = make(map[string]map[FieldType]*Field)
	}
	if f.Fields[name] == nil {
		f.Fields[name] = make(map[FieldType]*Field)
	}
	// Set index name to field.Indices
	field.Indices = f.Indices
	f.Fields[name][field.Type] = field
}
