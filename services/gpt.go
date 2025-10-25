package services

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"Oratio/models"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func ParseGPTResult(raw string) (*models.Session, error) {
	raw = strings.TrimSpace(raw)

	// Try to extract from ```json ... ```
	start := strings.Index(raw, "```json")
	end := strings.LastIndex(raw, "```")

	if start != -1 && end != -1 && end > start {
		raw = raw[start+7 : end] // skip "```json"
		raw = strings.TrimSpace(raw)
	}

	var result models.Session
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GPT(PaperBody string) (*models.Session, error) {
	ctx := context.Background()
	// The client gets the API key from the environment variable `OPENAI_API_KEY`.
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	
	client := openai.NewClient(apiKey)
	
	prompt := `You are an academic assistant. Given a research paper, do the following:

		1. Summarize the core message into a clear, spoken-style presentation of approximately 3 minutes.
		2. Generate 10–15 challenging audience questions that could be asked during the talk.
		3. Assign a random npc_id between 1 and 20 to each question.
		4. Return the result strictly in this JSON format (with no explanations or markdown):

		{
		"speech": "<speech_text_here>",
		"questions": [
			{ "npc_id": <int>, "text": "<question_1>" },
			{ "npc_id": <int>, "text": "<question_2>" }
			 ...
			{ "npc_id": <int>, "text": "<question_n>" }
		]
		}`

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini, // or openai.GPT4o, openai.GPT4Turbo
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: PaperBody,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	raw := resp.Choices[0].Message.Content // this is the big string
	parsed, err := ParseGPTResult(raw)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Speech:", parsed.Speech)

	return parsed, err
}
