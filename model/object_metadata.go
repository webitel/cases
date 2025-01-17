package model

type ObjectMetadatter interface {
	GetDefaultFields() []string
	GetAllFields() []string
	GetMainScopeName() string
	GetChildScopeNames() []string
	GetAllScopeNames() []string
}

type ObjectMetadata struct {
	fields               []string
	defFields            []string
	requiredObjClassName string
	childObjScopes       []string
	childMetadata        []ObjectMetadatter
}

func (o *ObjectMetadata) GetChildScopeNames() []string {
	return o.childObjScopes
}

func (o *ObjectMetadata) GetAllScopeNames() []string {
	return append(o.childObjScopes, o.requiredObjClassName)
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
func (o *ObjectMetadata) GetMainScopeName() string {
	return o.requiredObjClassName
}

type Field struct {
	Name    string
	Default bool
}

func NewObjectMetadata(requiredScope string, fields []*Field, childMetadata ...ObjectMetadatter) ObjectMetadatter {
	res := &ObjectMetadata{requiredObjClassName: requiredScope, childMetadata: childMetadata}
	for _, field := range fields {
		res.fields = append(res.fields, field.Name)
		if field.Default {
			res.defFields = append(res.defFields, field.Name)
		}
	}
	res.childMetadata = childMetadata
	uniqueMap := make(map[string]bool)

	for _, md := range childMetadata {
		if !uniqueMap[md.GetMainScopeName()] {
			uniqueMap[md.GetMainScopeName()] = true
			res.childObjScopes = append(res.childObjScopes, md.GetMainScopeName())
		}
	}

	return res
}
