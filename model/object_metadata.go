package model

type ObjectMetadatter interface {
	GetDefaultFields() []string
	GetAllFields() []string
	GetObjectName() string
}

type ObjectMetadata struct {
	fields    []string
	defFields []string
	objName   string
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
func (o *ObjectMetadata) GetObjectName() string {
	return o.objName
}

type Field struct {
	Name    string
	Default bool
}

func NewObjectMetadata(objName string, fields []*Field) *ObjectMetadata {
	res := &ObjectMetadata{objName: objName}
	for _, field := range fields {
		res.fields = append(res.fields, field.Name)
		if field.Default {
			res.defFields = append(res.defFields, field.Name)
		}
	}
	return res
}
