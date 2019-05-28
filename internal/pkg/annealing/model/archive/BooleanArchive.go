// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"math"
)

const actionsPerArchiveByte = 64

func BytesForManagementActions(managementActionNumber int) int {
	return int(math.Ceil(float64(managementActionNumber) / actionsPerArchiveByte))
}

type BooleanArchive struct {
	archiveBytes []byte
}

type archiveEntry struct {
	arrayIndex int
	byteOffset uint
	mask       byte
	value      bool
}

func New(size int) *BooleanArchive {
	newArchive := new(BooleanArchive)
	newArchive.archiveBytes = make([]byte, size)
	return newArchive
}

func (a *BooleanArchive) SetValue(entryIndex int, value bool) {
	entry := a.entry(entryIndex)

	if entry.value == value {
		return
	}

	arrayIndex := entry.arrayIndex

	if value {
		a.archiveBytes[arrayIndex] = a.archiveBytes[arrayIndex] + entry.mask
	} else {
		a.archiveBytes[arrayIndex] = a.archiveBytes[arrayIndex] - entry.mask
	}
}

func (a *BooleanArchive) Value(entryIndex int) bool {
	return a.entry(entryIndex).value
}

func (a *BooleanArchive) entry(entryIndex int) archiveEntry {
	index := entryIndex / actionsPerArchiveByte
	offset := uint(entryIndex % actionsPerArchiveByte)
	mask := byte(1 >> offset)
	value := a.archiveBytes[index]&mask > 0

	entry := archiveEntry{
		arrayIndex: index,
		byteOffset: offset,
		mask:       mask,
		value:      value,
	}

	return entry
}
