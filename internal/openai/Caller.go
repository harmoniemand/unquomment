package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func Call(promt string, ctx context.Context) (string, error) {

	oai_api_key := os.Getenv("OAI_API_KEY")

	oaiRequest := ChatcompletionRequest{
		Model:     ModelGpt35Turbo0125,
		MaxTokens: 300,
		Messages: []Message{
			{
				Role: MessageRoleSystem,
				Content: []Content{
					{
						Type: ContentTypeText,
						Text: promt,
					},
				},
			},
		},
	}

	reqJson, err := json.Marshal(oaiRequest)
	if err != nil {
		slog.Error("error marshaling json", "error", err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(reqJson))
	if err != nil {
		slog.Error("error creating request", "error", err)
		return "", err
	}

	reqWithCtx := req.WithContext(ctx)
	reqWithCtx.Header.Set("Authorization", "Bearer "+oai_api_key)
	reqWithCtx.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient

	res, err := client.Do(reqWithCtx)
	if err != nil {
		slog.Error("error doing request", "error", err)
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("error reading body", "error", err)
		return "", err
	}

	slog.Debug("body", "body", string(body))

	var chatcompletionResponse ChatcompletionResponse
	err = json.Unmarshal(body, &chatcompletionResponse)
	if err != nil {
		slog.Error("error unmarshaling json", "error", err)
		return "", err
	}

	text := chatcompletionResponse.Choices[0].Message.Content
	slog.Debug("text", "text", text)

	return text, nil
}
