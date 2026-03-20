package main

import (
	"fmt"
	"testing"

	"github.com/Rx-11/EDIS-A1/ai"
	"github.com/Rx-11/EDIS-A1/config"
)

func TestGeminiChat(t *testing.T) {
	config.Init() // Load environment variable for GEMINI_API_KEY

	if config.GetConfig().GeminiAPIKey == "" {
		t.Skip("Skipping test: GEMINI_API_KEY is not set in .env")
	}

	req := ai.ChatRequest{
		Messages: []ai.Message{
			{Role: "model", Content: "Give a 500 word summary of the following book"},
			{Role: "user", Content: "Book Title: The Great Gatsby\nBook Description: A novel set in the Jazz Age\nBook Author: F. Scott Fitzgerald\nBook ISBN: 978-0743273565"},
		},
	}

	resp, err := config.Gemini.Chat(req)
	if err != nil {
		t.Fatalf("Gemini.Chat() failed: %v", err)
	}

	if resp.Response == "" {
		t.Fatalf("Gemini.Chat() returned an empty response")
	}

	fmt.Println("------------- GEMINI SUMMARY -------------")
	fmt.Println(resp.Response)
	fmt.Println("------------------------------------------")
}
