package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/openai-chatGPT-server/config"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	. "github.com/openai-chatGPT-server/util"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type HistorySession struct {
	data sync.Map
	Cnf *config.ChatGPT
}


func (h *HistorySession) SetValue(key string, value interface{}) {
	msg := h.GetValue(key)

	newContent := append(msg, value.(openai.ChatCompletionMessage))
	h.data.Store(key, newContent)
}

func (h *HistorySession) GetValue(key string) []openai.ChatCompletionMessage {
	msg, ok := h.data.Load(key)
	if !ok {
		content := []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: h.Cnf.BotDesc},
		}
		msg =  content
	}
	historyContent := msg.([]openai.ChatCompletionMessage)
	return historyContent
}

func (h *HistorySession) Remove() {

}

func TurboStream(ctx *gin.Context, config *config.ChatGPT, msg []openai.ChatCompletionMessage) (content string, err error)  {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	var client = openai.NewClient(config.ApiKey)
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err != nil {
			Logger.Error(err)
			return "", err
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		clientConfig := openai.DefaultConfig(config.ApiKey)

		clientConfig.HTTPClient = &http.Client{
			Transport: transport,
		}
		client = openai.NewClientWithConfig(clientConfig)

	}
	request := openai.ChatCompletionRequest{
		Model: config.Model,
		Stream: true,
		Messages: msg,
	}
	Logger.Infof("request data:%+v", request)
	resp, err := client.CreateChatCompletionStream(
		context.Background(),
		request,
	)
	if err != nil {
		Logger.Error(err)
		return
	}

	for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			Logger.Info("Stream finished")
			break
		}

		if err != nil {
			Logger.Errorf("Stream error: %v\n", err)
			break
		}

		content += response.Choices[0].Delta.Content
		data, err := json.Marshal(response)
		if err != nil {
			Logger.Error(err)
			break
		}
		ctx.Writer.Write(data)
		ctx.Writer.Flush()
	}

	return content, nil
}

func Turbo(ctx *gin.Context, config *config.ChatGPT, msg []openai.ChatCompletionMessage) (content openai.ChatCompletionMessage, err error) {
	var client = openai.NewClient(config.ApiKey)
	if config.Proxy != "" {
		proxyUrl, err := url.Parse(config.Proxy)
		if err != nil {
			Logger.Error(err)
			return content, err
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		clientConfig := openai.DefaultConfig(config.ApiKey)

		clientConfig.HTTPClient = &http.Client{
			Transport: transport,
		}
		client = openai.NewClientWithConfig(clientConfig)

	}
	request := openai.ChatCompletionRequest{
		Model: config.Model,
		Messages: msg,
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		request,
	)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Infof("response:%+v", resp)

	return resp.Choices[0].Message, nil
}
