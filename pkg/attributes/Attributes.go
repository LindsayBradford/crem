// Copyright (c) 2019 Australian Rivers Institute.

// attributes package offers up a flexible approach to attaching attributes that then be used to avoid
// the brittle parameter problem.
package attributes

// Attributes is an array of name-value pairs.
type Attributes []NameValuePair

func (a *Attributes) Value(name string) interface{} {
	for _, attribute := range *a {
		if attribute.Name == name {
			return attribute.Value
		}
	}
	return nil
}

func (a Attributes) Entries(entries ...string) Attributes {
	slice := make(Attributes, len(entries))

	for _, attribute := range a {
		for entryIndex, entry := range entries {
			if attribute.Name == entry {
				slice[entryIndex] = attribute
			}
		}
	}

	return slice
}

func (a Attributes) Add(name string, value interface{}) Attributes {
	newEntry := NameValuePair{Name: name, Value: value}
	return append(a, newEntry)
}

func (a Attributes) Join(attributes Attributes) Attributes {
	return append(a, attributes...)
}

func (a Attributes) Rename(oldName string, newName string) Attributes {
	for index, attrib := range a {
		if attrib.Name == oldName {
			a[index].Name = newName
		}
	}
	return a
}

func (a Attributes) Remove(name string) Attributes {
	removeIndex := -1
	for index, attrib := range a {
		if attrib.Name == name {
			removeIndex = index
		}
	}

	frontAttributes := make(Attributes, 0)
	frontAttributes = a[:removeIndex]

	backAttributes := make(Attributes, 0)
	backAttributes = a[removeIndex+1:]

	return append(frontAttributes, backAttributes...)
}

func (a Attributes) Replace(name string, value interface{}) Attributes {
	for index, attrib := range a {
		if attrib.Name == name {
			a[index].Value = value
		}
	}
	return a
}

// NameValuePair is a struct allowing some Name as text to be associated with a matching Value of any type.
type NameValuePair struct {
	Name  string
	Value interface{}
}
