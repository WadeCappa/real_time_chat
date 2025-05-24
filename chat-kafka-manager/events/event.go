package events

type Event interface {
	Visit(EventVisitor)
}
