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

func main () {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

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

	viper.ReadInConfig()
	apiKey := viper.GetString("api_key")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	
	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text("How are you today?"))
	if err != nil {
		fmt.Println(err)
	}

	printResponse(resp)	
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}