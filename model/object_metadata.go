package model

type ObjectMetadatter interface {
	GetDefaultFields() []string
	GetAllFields() []string
}

type ObjectMetadata struct {
	fields    []string
	defFields []string
}

func (o *ObjectMetadata) GetAllFields() []string {
	res := make([]string, len(o.fields))
	copy(res, o.fields)
	return res
}

func (o *ObjectMetadata) GetDefaultFields() []string {
	res := make([]string, len(o.defFields))
	copy(res, o.defFields)
	return res
}

type Field struct {
	Name    string
	Default bool
}

func NewObjectMetadata(fields []*Field) *ObjectMetadata {
	res := &ObjectMetadata{}
	for _, field := range fields {
		res.fields = append(res.fields, field.Name)
		if field.Default {
			res.defFields = append(res.defFields, field.Name)
		}
	}
	return res
}
