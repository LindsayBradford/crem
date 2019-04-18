// Copyright (c) 2019 Australian Rivers Institute.

package attributes

type ContainedAttributes struct {
	attributes Attributes
}

func (ca *ContainedAttributes) Attribute(name string) interface{} {
	return ca.attributes.Value(name)
}

func (ca *ContainedAttributes) AddAttribute(name string, value interface{}) {
	newEntry := NameValuePair{Name: name, Value: value}
	ca.attributes = append(ca.attributes, newEntry)
}
