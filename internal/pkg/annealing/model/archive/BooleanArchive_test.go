// Copyright (c) 2019 Australian Rivers Institute.

package archive

import (
	"testing"

	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestBankSedimentContribution_BytesForManagementActions(t *testing.T) {
	g := NewGomegaWithT(t)

	var numberOfManagementActions int
	var expectedBytes int
	var actualBytes int

	numberOfManagementActions = 0
	expectedBytes = 0

	actualBytes = BytesForManagementActions(numberOfManagementActions)
	g.Expect(actualBytes).To(BeNumerically(equalTo, expectedBytes))

	numberOfManagementActions = 1
	expectedBytes = 1

	actualBytes = BytesForManagementActions(numberOfManagementActions)
	g.Expect(actualBytes).To(BeNumerically(equalTo, expectedBytes))

	numberOfManagementActions = 64
	expectedBytes = 1

	actualBytes = BytesForManagementActions(numberOfManagementActions)
	g.Expect(actualBytes).To(BeNumerically(equalTo, expectedBytes))

	numberOfManagementActions = 65
	expectedBytes = 2

	actualBytes = BytesForManagementActions(numberOfManagementActions)
	g.Expect(actualBytes).To(BeNumerically(equalTo, expectedBytes))

	numberOfManagementActions = 128
	expectedBytes = 2

	actualBytes = BytesForManagementActions(numberOfManagementActions)
	g.Expect(actualBytes).To(BeNumerically(equalTo, expectedBytes))

	numberOfManagementActions = 129
	expectedBytes = 3

	actualBytes = BytesForManagementActions(numberOfManagementActions)
	g.Expect(actualBytes).To(BeNumerically(equalTo, expectedBytes))
}
