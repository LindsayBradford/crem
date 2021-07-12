// Copyright (c) 2019 Australian Rivers Institute.

package attributes

type Interface interface {
	HasAttribute(name string) bool
	Attribute(name string) interface{}
	AllAttributes() Attributes

	AddAttribute(name string, value interface{})
	RenameAttribute(currentName string, newName string)
	ReplaceAttribute(name string, value interface{})
	RemoveAttribute(name string)

	JoiningAttributes(newAttributes Attributes)
}

type ContainedAttributes struct {
	Attribs Attributes `json:"Attributes"`
}

var _ Interface = new(ContainedAttributes)

func (ca *ContainedAttributes) Attribute(name string) interface{} {
	return ca.Attribs.Value(name)
}

func (ca *ContainedAttributes) HasAttribute(name string) bool {
	return ca.Attribs.Value(name) != nil
}

func (ca *ContainedAttributes) AddAttribute(name string, value interface{}) {
	ca.Attribs = ca.Attribs.Add(name, value)
}

func (ca *ContainedAttributes) RenameAttribute(currentName string, newName string) {
	ca.Attribs = ca.Attribs.Rename(currentName, newName)
}

func (ca *ContainedAttributes) ReplaceAttribute(name string, value interface{}) {
	if ca.HasAttribute(name) {
		ca.Attribs = ca.Attribs.Replace(name, value)
	} else {
		ca.AddAttribute(name, value)
	}
}

func (ca *ContainedAttributes) RemoveAttribute(name string) {
	if !ca.HasAttribute(name) {
		return
	}
	ca.Attribs = ca.Attribs.Remove(name)
}

func (ca *ContainedAttributes) JoiningAttributes(newAttributes Attributes) {
	ca.Attribs = ca.Attribs.Join(newAttributes)
}

func (ca *ContainedAttributes) AttributesNamed(entries ...string) Attributes {
	return ca.Attribs.Entries(entries...)
}

func (ca *ContainedAttributes) AllAttributes() Attributes {
	return ca.Attribs
}
