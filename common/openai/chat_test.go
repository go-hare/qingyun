package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestClient_CreateChatCompletion(t *testing.T) {
	conf := DefaultConfig("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")
	hc := http.Client{Timeout: 30 * time.Second}
	hc.Transport = &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:7890")
		},
	}
	conf.HTTPClient = &hc
	client := NewClientWithConfig(conf)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		ChatCompletionRequest{
			Model: GPT3Dot5Turbo,
			Messages: []ChatCompletionMessage{
				{
					Role:    ChatMessageRoleUser,
					Content: "中问",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}

func TestClient_CreateChatCompletionStream(t *testing.T) {
	conf := DefaultConfig("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")
	hc := http.Client{Timeout: 30 * time.Second}
	hc.Transport = &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:7890")
		},
	}
	conf.HTTPClient = &hc
	client := NewClientWithConfig(conf)
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		ChatCompletionRequest{
			Model:     GPT3Dot5Turbo,
			MaxTokens: 20,
			Messages: []ChatCompletionMessage{
				{
					Role:    ChatMessageRoleUser,
					Content: "你好22222222",
				},
			},
			Stream: true,
		},
	)
	defer stream.Close()
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}
		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
