package mvp

type ViewEventType int

func (vet *ViewEventType) String() string {
	switch *vet {
	case GenerationRequested:
		return "GenerationRequested"
	case FocusGained:
		return "FocusGained"

	default:
		return "<undefined>"
	}
}

const (
	GenerationRequested ViewEventType = iota
	FocusGained
)

type ViewEvent struct {
	View    View
	Type    ViewEventType
	Context interface{}
}

type View interface {
	Id() string
	AddObserver(o ViewObserver)
	Show()
	SetMessage(string)
}

type ViewObserver interface {
	EventRaised(ve ViewEvent)
}
