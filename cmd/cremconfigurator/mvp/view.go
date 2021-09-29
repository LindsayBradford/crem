package mvp

type ViewEvent int

const (
	GenerationRequested ViewEvent = iota
)

type View interface {
	AddObserver(o ViewObserver)
	Show()
}

type ViewObserver interface {
	EventRaised(ve ViewEvent)
}
