package dream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Dream struct {
	ID        string
	Title     string
	CreatedAt time.Time
	RawText   string
}

func NewDream(id, title, rawText string, createdAt time.Time) (*Dream, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id is blank")
	}

	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title is blank")
	}

	if strings.TrimSpace(rawText) == "" {
		return nil, fmt.Errorf("content is blank")
	}

	if createdAt.IsZero() {
		return nil, fmt.Errorf("created at not provided")
	}

	return &Dream{
		ID: id, Title: title, RawText: rawText, CreatedAt: createdAt,
	}, nil
}

type DreamAnalysis struct {
	ID        int
	DreamID   string
	Summary   string
	Themes    []string
	Emotions  []string
	Locations []string
	People    []string
	Symbols   []string
}

type DreamRepository interface {
	CreateDream(d Dream) error
	ListDreams() ([]Dream, error)
	GetDream(id string) (Dream, error)
	DeleteDream(id string) error

	SaveDreamAnalysis(da DreamAnalysis) error
	GetDreamAnalysis(dreamID string) (*DreamAnalysis, error)
}

type DreamAnalyser interface {
	Analyse(ctx context.Context, d Dream) (*DreamAnalysis, error)
}

type FakeAnalyser struct{}

func (a *FakeAnalyser) Analyse(ctx context.Context, d Dream) (*DreamAnalysis, error) {
	return &DreamAnalysis{
		DreamID:   d.ID,
		Summary:   "A test analysis",
		Themes:    []string{"testing", "romance"},
		Emotions:  []string{"misery"},
		Locations: []string{"Belfast", "Derry"},
		People:    []string{"Eamonn", "Amy"},
		Symbols:   []string{"Trophy", "Red Rose"},
	}, nil
}

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

func (a *LocalLLMAnalyser) Analyse(ctx context.Context, d Dream) (*DreamAnalysis, error) {
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
					- notable objects or symbols
					Be extra careful to assign the correct actions to each person in the dream.
					If there are multiple people in the dream, do not mix up which person performed which action or 
					which person had which thing happen to them. This is crucial.

					If a category has no information, return an empty array.

					Return ONLY the JSON object described by the response schema.
					Do not use Markdown.
					Do not provide an explanation. 
					There should be zero use of backticks in your response.`,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Analyse this dream: \n\n%s", d.RawText),
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

	data, _ := json.MarshalIndent(reqBody, "", "  ")
	fmt.Println(string(data))

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var chatResponse chatResponse

	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return nil, err
	}

	content := chatResponse.Choices[0].Message.Content
	fmt.Printf("%q", chatResponse)

	var result localLLMAnalysisResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, err
	}

	return &DreamAnalysis{
		DreamID:   d.ID,
		Summary:   result.Summary,
		Themes:    result.Themes,
		Emotions:  result.Emotions,
		Locations: result.Locations,
		People:    result.People,
		Symbols:   result.Symbols,
	}, nil
}
