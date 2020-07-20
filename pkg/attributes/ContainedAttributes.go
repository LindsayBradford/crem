// Copyright (c) 2019 Australian Rivers Institute.

package attributes

type ContainedAttributes struct {
	attributes Attributes
}

func (ca *ContainedAttributes) Attribute(name string) interface{} {
	return ca.attributes.Value(name)
}

func (ca *ContainedAttributes) HasAttribute(name string) bool {
	return ca.attributes.Value(name) != nil
}

func (ca *ContainedAttributes) AddAttribute(name string, value interface{}) {
	ca.attributes = ca.attributes.Add(name, value)
}

func (ca *ContainedAttributes) RenameAttribute(currentName string, newName string) {
	ca.attributes = ca.attributes.Rename(currentName, newName)
}

func (ca *ContainedAttributes) ReplaceAttribute(name string, value interface{}) {
	if ca.HasAttribute(name) {
		ca.attributes = ca.attributes.Replace(name, value)
	} else {
		ca.AddAttribute(name, value)
	}
}

func (ca *ContainedAttributes) RemoveAttribute(name string) {
	if !ca.HasAttribute(name) {
		return
	}
	ca.attributes = ca.attributes.Remove(name)
}

func (ca *ContainedAttributes) JoiningAttributes(newAttributes Attributes) {
	ca.attributes = ca.attributes.Join(newAttributes)
}

func (ca *ContainedAttributes) AttributesNamed(entries ...string) Attributes {
	return ca.attributes.Entries(entries...)
}

func (ca *ContainedAttributes) AllAttributes() Attributes {
	return ca.attributes
}
