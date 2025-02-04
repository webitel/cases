package model

type ObjectMetadatter interface {
	GetDefaultFields() []string
	GetAllFields() []string
	GetMainScopeName() string
	GetParentScopeName() string
	GetChildScopeNames() []string
	GetAllScopeNames() []string
}

type ObjectMetadata struct {
	fields             []string
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
	return o.mainObjClassName
}

type Field struct {
	Name    string
	Default bool
}

func NewObjectMetadata(mainScope string, parentScope string, fields []*Field, childMetadata ...ObjectMetadatter) ObjectMetadatter {
	res := &ObjectMetadata{mainObjClassName: mainScope, parentObjClassName: parentScope, childMetadata: childMetadata}
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

// SetAllFieldsToTrue returns a copy of ObjectMetadata with all fields set as default (true)
func SetAllFieldsToTrue(o ObjectMetadata) ObjectMetadata {
	// Create a deep copy of ObjectMetadata
	newMetadata := ObjectMetadata{
		mainObjClassName:   o.mainObjClassName,
		parentObjClassName: o.parentObjClassName,
		childObjScopes:     append([]string{}, o.childObjScopes...),          // Copy slice
		childMetadata:      append([]ObjectMetadatter{}, o.childMetadata...), // Copy slice
	}

	// Copy and modify fields
	newMetadata.fields = append([]string{}, o.fields...)    // Copy field names
	newMetadata.defFields = append([]string{}, o.fields...) // Set all fields as default

	return newMetadata
}
