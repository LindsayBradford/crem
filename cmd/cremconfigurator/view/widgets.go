package view

import (
	"fyne.io/fyne/v2/widget"
	"github.com/LindsayBradford/crem/cmd/cremconfigurator/mvp"
)

type cremEntry struct {
	widget.Entry
	Label string
}

func NewCremEntry(label string) *cremEntry {
	newEntry := cremEntry{
		Entry: *widget.NewEntry(),
		Label: label,
	}

	return &newEntry
}

func (ce cremEntry) FocusGained() {
	ce.Entry.FocusGained()
	view.raiseEvent(mvp.ViewEvent{
		Type:    mvp.FocusGained,
		Context: ce.Label,
	},
	)
}
