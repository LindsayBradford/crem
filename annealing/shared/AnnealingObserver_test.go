// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestAnnealingEventType_String(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(InvalidEvent.String()).To(Equal("InvalidEvent"))

	tooSmallEventType := InvalidEvent - 1
	g.Expect(tooSmallEventType.String()).To(Equal("InvalidEvent"))

	g.Expect(Note.String()).To(Equal("Note"))

	tooLargeEventType := Note + 1
	g.Expect(tooLargeEventType.String()).To(Equal("InvalidEvent"))
}
