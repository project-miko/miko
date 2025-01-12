package chatgptapi

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/project-miko/miko/conf"
	"github.com/sashabaranov/go-openai"
)

var (
	defaultTimeOut          = 180
	apiKey                  = ""
	maxTokens       int     = 4000
	temperature     float32 = 0.7
	modelStr                = openai.GPT4o
	presencePenalty float32
)

func InitChatGPT() error {
	apiKey = conf.GetConfigString("chatgpt", "api_key")
	if len(apiKey) == 0 {
		return fmt.Errorf("chatgpt.api_key not config")
	}
	var err error
	_maxTokens, err := conf.GetConfigInt("chatgpt", "max_tokens")
	if err != nil {
		return err
	}
	maxTokens = int(_maxTokens)
	_str := conf.GetConfigString("chatgpt", "temperature")
	_temperature, err := strconv.ParseFloat(_str, 32)
	if err != nil {
		return err
	}
	temperature = float32(_temperature)
	ppStr := conf.GetConfigString("chatgpt", "presence_penalty")
	_presencePenalty, err := strconv.ParseFloat(ppStr, 32)
	if err != nil {
		return err
	}
	presencePenalty = float32(_presencePenalty)

	return nil
}

func SendChatGPTRequest(roleSystemContent, roleUserContent string) (string, error) {

	// ChatCompletion error: error, status code: 400, message:
	//   This model's maximum context length is 4097 tokens. However, your messages resulted in 15105 tokens.
	//   Please reduce the length of the messages.
	if utf8.RuneCountInString(roleUserContent) > maxTokens {
		strF := []rune(roleUserContent)
		roleUserContent = string(strF[:maxTokens])
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeOut)*time.Second)
	defer cancel()

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:           modelStr,
			Temperature:     temperature,
			PresencePenalty: presencePenalty,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: roleSystemContent,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: roleUserContent,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if resp.Choices == nil || len(resp.Choices) <= 0 {
		return "", fmt.Errorf("choices is empty, response %v", resp)
	}

	rawContent := resp.Choices[0].Message.Content

	return rawContent, nil
}
