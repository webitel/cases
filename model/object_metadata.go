package model

type ObjectMetadatter interface {
	GetDefaultFields() []string
	GetAllFields() []string
	GetMainScopeName() string
	GetParentScopeName() string
	GetChildScopeNames() []string
	GetAllScopeNames() []string
	SetAllFieldsToTrue() *ObjectMetadata
}
type Field struct {
	Name    string
	Default bool
}
type ObjectMetadata struct {
	fields             []*Field
	defFields          []string
	mainObjClassName   string
	parentObjClassName string
	childObjScopes     []string
	childMetadata      []ObjectMetadatter
}

func (o *ObjectMetadata) GetParentScopeName() string {
	return o.parentObjClassName
}

func (o *ObjectMetadata) GetChildScopeNames() []string {
	return o.childObjScopes
}

func (o *ObjectMetadata) GetAllScopeNames() []string {
	return append(o.childObjScopes, o.mainObjClassName, o.parentObjClassName)
}

func (o *ObjectMetadata) GetAllFields() []string {
	var res []string
	for _, field := range o.fields {
		res = append(res, field.Name)
	}
	return res
}

func (o *ObjectMetadata) GetDefaultFields() []string {
	var res []string
	for _, field := range o.defFields {
		res = append(res, field)
	}
	return res
}

func (o *ObjectMetadata) GetMainScopeName() string {
	return o.mainObjClassName
}

// NewObjectMetadata creates a new ObjectMetadata instance
func NewObjectMetadata(mainScope string, parentScope string, fields []*Field, childMetadata ...ObjectMetadatter) ObjectMetadatter {
	res := &ObjectMetadata{mainObjClassName: mainScope, parentObjClassName: parentScope, childMetadata: childMetadata}
	for _, field := range fields {
		res.fields = append(res.fields, field)
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

// SetAllFieldsToTrue sets all fields' Default property to true
func (o *ObjectMetadata) SetAllFieldsToTrue() *ObjectMetadata {
	// Loop through all fields and set Default to true
	for _, field := range o.fields {
		field.Default = true // Set the Default to true for each field
	}

	// Return the updated ObjectMetadata
	return o
}
