package events

type NewChatMessageEvent struct {
	Event

	Content   string
	UserId    int64
	ChannelId int64
}

func (e NewChatMessageEvent) Visit(v EventVisitor) error {
	return v.VisitNewChatMessageEvent(e)
}
