package openai

//	{
//		"id": "chatcmpl-8ijAsaCzlw8TTOPxog6QJJu8KAEwJ",
//		"object": "chat.completion",
//		"created": 1705670774,
//		"model": "gpt-4-1106-vision-preview",
//		"usage": {
//			"prompt_tokens": 814,
//			"completion_tokens": 70,
//			"total_tokens": 884
//		},
//		"choices": [
//			{
//				"message": {
//					"role": "assistant",
//					"content": "Du stehst vor einem Fußgängerüberweg mit Ampel. Die Ampel zeigt Grün für Fußgänger. Einige Personen überqueren die Straße. Auf der rechten Seite sind Fahrradwege. Achte auf vorbeifahrende Fahrräder, wenn du die Straße überquerst."
//				},
//				"finish_reason": "stop",
//				"index": 0
//			}
//		]
//	}
type ChoiceMessage struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

type Choice struct {
	Message      ChoiceMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
	Index        int           `json:"index"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatcompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}
