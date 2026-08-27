package provider

import (
	"context"
	"net/http"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type OpenRouterProvider struct {
	Client 	*openai.Client
	model 	string
	logger 	*zap.SugaredLogger
}

func NewOpenRouterProvider (apiKey, baseURL, model string, logger *zap.SugaredLogger) *OpenRouterProvider {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: false,
		},
	}
	return &OpenRouterProvider{
		Client: openai.NewClientWithConfig(config),
		model: model,
		logger: logger,
	}
}

func (p *OpenRouterProvider) GeneratePlan(ctx context.Context, systemPrompt , userPrompt string) (string, error) {
	p.logger.Infof("Sending request to model: %s", p.model)
	resp, err := p.Client.CreateChatCompletion( 
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			MaxTokens: 3000,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
			 	Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			Messages: []openai.ChatCompletionMessage {
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: userPrompt},
			},
		},
	)
	if err != nil {
		p.logger.Errorf("OpenRouter API error %v",err)
		return "",err
		
	}
	
	return resp.Choices[0].Message.Content, nil 
}
