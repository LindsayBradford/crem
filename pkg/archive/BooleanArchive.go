// Copyright (c) 2019 Australian Rivers Institute.

// Package archive offers support for the compact storage and retrieval of state in memory.
package archive

import (
	"github.com/pkg/errors"
	"math"
	"sync"
)

const entriesPerArchiveEntry = 64

// BooleanArchive offers a compact storage and retrieval mechanism for a fixed size set of boolean data.
type BooleanArchive struct {
	size int

	cacheMutex  sync.Mutex
	detailCache entryDetail

	archiveArray []uint64
}

type entryDetail struct {
	arrayIndex int
	byteOffset uint
	mask       uint64

	value bool
}

// BooleanArchive offers a compact storage and retrieval mechanism for a fixed size set of boolean data.
func New(size int) *BooleanArchive {
	newArchive := new(BooleanArchive)

	newArchive.size = size
	newArchive.archiveArray = make([]uint64, archiveSize(size))

	return newArchive
}

// Len returns the length of entries stored in the archive.
func (a *BooleanArchive) Len() int {
	return a.size
}

// ArchiveLen returns the number of uint64 integers used to store archived boolean data.
func (a *BooleanArchive) ArchiveLen() int {
	return len(a.archiveArray)
}

// SetValue stores the supplied boolean value in the archive at entryIndex.
func (a *BooleanArchive) SetValue(entryIndex int, value bool) {
	if entryIndex >= a.size {
		outOfBoundsError := errors.New("index out of range")
		panic(outOfBoundsError)
	}

	a.cacheMutex.Lock()
	defer a.cacheMutex.Unlock()

	entry := a.deriveDetail(entryIndex)

	if entry.value == value {
		return
	}

	arrayIndex := entry.arrayIndex

	if value {
		a.archiveArray[arrayIndex] = a.archiveArray[arrayIndex] + entry.mask
	} else {
		a.archiveArray[arrayIndex] = a.archiveArray[arrayIndex] - entry.mask
	}
}

// Value retrieves the boolean value stored at the requested entryIndex of the archive.
func (a *BooleanArchive) Value(entryIndex int) bool {
	if entryIndex >= a.size {
		outOfBoundsError := errors.New("index out of range")
		panic(outOfBoundsError)
	}

	a.cacheMutex.Lock()
	defer a.cacheMutex.Unlock()

	return a.deriveDetail(entryIndex).value
}

func (a *BooleanArchive) deriveDetail(entryIndex int) *entryDetail {
	a.detailCache.arrayIndex = entryIndex / entriesPerArchiveEntry
	a.detailCache.byteOffset = uint(entryIndex % entriesPerArchiveEntry)
	a.detailCache.mask = uint64(1 << a.detailCache.byteOffset)
	a.detailCache.value = a.archiveArray[a.detailCache.arrayIndex]&a.detailCache.mask > 0

	return &a.detailCache
}

func archiveSize(entriesNeeded int) int {
	return int(math.Ceil(float64(entriesNeeded) / entriesPerArchiveEntry))
}

func (a *BooleanArchive) IsEquivalentTo(b *BooleanArchive) bool {
	if a.size != b.size {
		return false
		return false
	}

	for index := range a.archiveArray {
		if a.archiveArray[index] != b.archiveArray[index] {
			return false
		}
	}
	return true
}
