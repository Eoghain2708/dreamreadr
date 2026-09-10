package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
)

type LocalLLMAnalyser struct {
	Client  *http.Client
	BaseURL string
}

func NewLocalLLMAnalyser(c *http.Client, url string) (*LocalLLMAnalyser, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("url is blank")
	}

	return &LocalLLMAnalyser{Client: c, BaseURL: url}, nil
}

// expected JSON response from LLM
type localLLMAnalysisResponse struct {
	Summary   string   `json:"summary"`
	Themes    []string `json:"themes"`
	Emotions  []string `json:"emotions"`
	Locations []string `json:"locations"`
	People    []string `json:"people"`
	Symbols   []string `json:"symbols"`
}

type chatRequest struct {
	Messages       []message      `json:"messages"`
	ResponseFormat responseFormat `json:"response_format"`
}

type responseFormat struct {
	Type       string         `json:"type"`
	JSONSchema map[string]any `json:"schema,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (a *LocalLLMAnalyser) Analyse(ctx context.Context, d dream.Dream) (*dream.DreamAnalysis, error) {
	reqBody := chatRequest{
		Messages: []message{
			{
				Role: "system",
				Content: `You are a dream journal information extraction system.

					Your task is to extract details that are explicitly present
					in the dream narrative.

					You are NOT a psychologist.
					You are NOT a therapist.
					You must NOT interpret the dream.
					You must NOT diagnose the dreamer.
					You must NOT speculate about the dreamer's mental state.
					You must NOT assign psychological meanings to symbols.

					For example:

					If the dream says:
					"I was terrified while walking through a hospital."

					Then:
					emotions = ["fear"]
					locations = ["hospital"]

					Do NOT infer:
					"the hospital represents anxiety about health."

					Only record information that is directly supported
					by the narrative.

					Extract:
					- a concise factual summary of what happened
					- prominent narrative themes explicitly present
					- emotions experienced or explicitly described
					- locations
					- people
					- notable objects or symbols, i.e "red rose", "snooker cue", things that the user could learn to recognise when 
					learning to lucid dream.
					Be extra careful to assign the correct actions to each person in the dream.
					If there are multiple people in the dream, do not mix up which person performed which action or 
					which person had which thing happen to them. This is crucial.

					If a category has no information, return an empty array.
					If there is no storyline to a dream, do not make one up. For example, if given the dream "This is a test dream",
					do not make up a story where there is none.

					Return ONLY the JSON object described by the response schema.
					Do not use Markdown.
					Do not provide an explanation. 
					There should be zero use of backticks in your response.`,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Extract the information from this dream: \n\n%s", d.RawText),
			},
		},

		ResponseFormat: responseFormat{
			Type: "json_object",
			JSONSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"summary": map[string]any{
						"type": "string",
					},
					"themes": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
					"emotions": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
					"locations": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
					"people": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
					"symbols": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
				"required": []string{
					"summary",
					"themes",
					"emotions",
					"locations",
					"people",
					"symbols",
				},
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.BaseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var chatResponse chatResponse

	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return nil, err
	}

	content := chatResponse.Choices[0].Message.Content

	var result localLLMAnalysisResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, err
	}

	return &dream.DreamAnalysis{
		DreamID:   d.ID,
		Summary:   result.Summary,
		Themes:    result.Themes,
		Emotions:  result.Emotions,
		Locations: result.Locations,
		People:    result.People,
		Symbols:   result.Symbols,
	}, nil
}
