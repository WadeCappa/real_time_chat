package events

type NewChatMessageEvent struct {
	Content   string
	UserId    int64
	ChannelId int64
}

func (e NewChatMessageEvent) Visit(v EventVisitor) {
	v.visitNewChatMessageEvent(e)
}
