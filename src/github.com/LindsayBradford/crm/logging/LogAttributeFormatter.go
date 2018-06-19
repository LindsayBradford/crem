// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

const NULL_FORMAT_MESSAGE = "No formatter specified. Using the NullFormatter."

type NullFormatter struct {}

func (this *NullFormatter) Initialise() {}

func (this *NullFormatter) Format(attributes LogAttributes) string {
	return NULL_FORMAT_MESSAGE
}

const DEFAULT_MESSAGE_LABEL = "message"

type MessageFormatter struct {
	messageLabel string
}

func (this *MessageFormatter) Initialise() {
	this.messageLabel = DEFAULT_MESSAGE_LABEL
}

func (this *MessageFormatter) Format(attributes LogAttributes) string {
	for _, attribute := range attributes {
		if (attribute.Name == this.messageLabel) {
			return attribute.Value.(string)
		}
	}
	return ""
}
