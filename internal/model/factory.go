package model

func NewMessageFactory() MessageFactory {
	return MessageFactory{}
}

type MessageFactory struct {
}

func (m MessageFactory) CreateWithContent(content string) Message {
	return Message{Content: content}
}
