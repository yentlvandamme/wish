package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Pass a string value to act as your prompt")
		os.Exit(1)
	}

	prompt := os.Args[1]

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	createConfigIfNotExist()

	viper.ReadInConfig()
	apiKey := viper.GetString("api_key")

	var llmHandler LLMHandler
	geminiHandler, err := NewGeminiHandler(apiKey)
	if err != nil {
		fmt.Println(err)
	}

	llmHandler = geminiHandler
	parts, err := llmHandler.Ask(prompt)
	if err != nil {
		fmt.Println(err)
	}

	for _, part := range parts {
		fmt.Println(part)
	}
}

type LLMHandler interface {
	Ask(prompt string) ([]string, error)
}

type GeminiHandler struct {
	apiKey  string
	context context.Context
	client  *genai.Client
	model   *genai.GenerativeModel
}

func NewGeminiHandler(apiKey string) (*GeminiHandler, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))

	if err != nil {
		fmt.Println("Could not create Gemini client")
		return nil, err
	}

	return &GeminiHandler{
		apiKey:  apiKey,
		context: ctx,
		client:  client,
		model:   client.GenerativeModel("gemini-1.5-flash"),
	}, nil
}

func (handler GeminiHandler) Ask(prompt string) ([]string, error) {
	generatedContent, err := handler.model.GenerateContent(handler.context, genai.Text(prompt))
	responseParts := make([]string, 0)

	if err != nil {
		fmt.Println("Error getting the response")
		return responseParts, err
	}

	for _, candidate := range generatedContent.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					responseParts = append(responseParts, string(txt))
				} else {
					fmt.Println("Generated content is not of type genai.Text")
				}
			}
		}
	}

	return responseParts, nil
}

func createConfigIfNotExist() {
	if _, err := os.Stat("./config.yaml"); os.IsNotExist(err) {
		_, err := os.OpenFile("./config.yaml", os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("API key: ")
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Invalid input", err)
		}

		viper.Set("api_key", strings.TrimSpace(apiKey))
		if err := viper.WriteConfig(); err != nil {
			fmt.Println(err)
		}
	}
}
