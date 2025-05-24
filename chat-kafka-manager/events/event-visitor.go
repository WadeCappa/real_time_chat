package events

type EventVisitor interface {
	visitNewChatMessageEvent(e NewChatMessageEvent)
}
