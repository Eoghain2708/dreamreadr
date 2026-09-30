package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"google.golang.org/genai"
)

var schema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"summary": {
			Type: genai.TypeString,
		},
		"themes": {
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
		"emotions": {
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
		"symbols": {
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
		"locations": {
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
		"people": {
			Type:  genai.TypeArray,
			Items: &genai.Schema{Type: genai.TypeString},
		},
	},
}

type GeminiLLMAnalyser struct {
	Client *genai.Client
	Model  string
}

func NewGeminiLLMAnalyser(model string) (*GeminiLLMAnalyser, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &GeminiLLMAnalyser{Client: client, Model: model}, nil
}

func (a *GeminiLLMAnalyser) Analyse(ctx context.Context, d dream.Dream) (*dream.DreamAnalysis, error) {
	parts := []*genai.Part{
		{Text: fmt.Sprintf(`You are a dream journal information extraction system.

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
					There should be zero use of backticks in your response.
					Extract the information from this dream: \n\n%s`, d.RawText)},
	}

	config := genai.GenerateContentConfig{ResponseMIMEType: "application/json", ResponseSchema: schema}

	res, err := a.Client.Models.GenerateContent(ctx, a.Model, []*genai.Content{{Parts: parts}}, &config)
	if err != nil {
		return nil, err
	}

	text := res.Text()

	var da *dream.DreamAnalysis

	if err := json.Unmarshal([]byte(text), &da); err != nil {
		return nil, err
	}

	da.DreamID = d.ID
	return da, nil
}
