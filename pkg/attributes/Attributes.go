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
	slice := make(Attributes, 0)

	for _, attribute := range a {
		for _, entry := range entries {
			if attribute.Name == entry {
				slice = append(slice, attribute)
			}
		}
	}

	return slice
}

// NameValuePair is a struct allowing some Name as text to be associated with a matching Value of any type.
type NameValuePair struct {
	Name  string
	Value interface{}
}
