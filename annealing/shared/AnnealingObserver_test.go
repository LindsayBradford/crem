// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestAnnealingEventType_String(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(INVALID_EVENT.String()).To(Equal("INVALID_EVENT"))

	tooSmallEventType := INVALID_EVENT - 1
	g.Expect(tooSmallEventType.String()).To(Equal("INVALID_EVENT"))

	g.Expect(NOTE.String()).To(Equal("NOTE"))

	tooLargeEventType := NOTE + 1
	g.Expect(tooLargeEventType.String()).To(Equal("INVALID_EVENT"))
}
