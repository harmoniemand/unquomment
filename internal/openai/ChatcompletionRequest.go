package openai

import (
	"encoding/json"
	"log/slog"
)

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImageUrl ContentType = "image_url"
)

func (c ContentType) MarshalJSON() ([]byte, error) {
	str := string(c)
	return json.Marshal(str)
}

func (c *ContentType) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	switch str {
	case "text":
		*c = ContentTypeText
	case "image_url":
		*c = ContentTypeImageUrl
	default:
		slog.Error("unknown ContentType", "ContentType", str)
		return nil
	}

	return nil
}

// func (c ContentType) UnmarshalJSON(b []byte,

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleSystem    MessageRole = "system"
	MessageRoleAssistant MessageRole = "assistant"
)

func (r MessageRole) MarshalJSON() ([]byte, error) {
	str := string(r)
	return json.Marshal(str)
}

func (r *MessageRole) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	switch str {
	case "user":
		*r = MessageRoleUser
	case "system":
		*r = MessageRoleSystem
	case "assistant":
		*r = MessageRoleAssistant
	default:
		slog.Error("unknown MessageRole", "MessageRole", str)
		return nil
	}

	return nil
}

type Model string

const (
	ModelGpt4VisionPreview Model = "gpt-4-vision-preview"
	ModelGpt4TurboPreview  Model = "gpt-4-0125-preview"
	ModelGpt35Turbo0125    Model = "gpt-3.5-turbo-0125"
)

func (m Model) MarshalJSON() ([]byte, error) {
	str := string(m)
	return json.Marshal(str)
}

type ImageUrl struct {
	Url string `json:"url"`
}

type Content struct {
	Type     ContentType `json:"type"`
	Text     string      `json:"text,omitempty"`
	ImageUrl *ImageUrl   `json:"image_url,omitempty"`
}

type Message struct {
	Role    MessageRole `json:"role"`
	Content []Content   `json:"content"`
}

type ChatcompletionRequest struct {
	Model     Model     `json:"model"`
	MaxTokens int       `json:"max_tokens" default:"300"`
	Messages  []Message `json:"messages"`
}
