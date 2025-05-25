package events

type EventVisitor interface {
	VisitNewChatMessageEvent(e NewChatMessageEvent) error
}
