// Copyright (c) 2020 Australian Rivers Institute.

package variable

import "sort"

const KeyNotFound = -1

type DecisionVariableMap map[string]DecisionVariable

func (m DecisionVariableMap) SortedKeys() (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func (m DecisionVariableMap) SortedKeyIndex(key string) int {
	sortedKeys := m.SortedKeys()
	for index, currentKey := range sortedKeys {
		if currentKey == key {
			return index
		}
	}
	return KeyNotFound
}
