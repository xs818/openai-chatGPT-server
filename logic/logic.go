package logic

import (
	"github.com/openai-chatGPT-server/config"
	. "github.com/openai-chatGPT-server/util"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
)

type Chat interface  {
	ChatRoomStream(ctx *gin.Context)
	ChatRoom(ctx *gin.Context)
}

type chat struct {
	cnf *config.Configuration
	history *HistorySession
}

func NewChat(cnf *config.Configuration) Chat {
	history := &HistorySession{Cnf: &cnf.ChatGPT}
	return &chat{
		cnf: cnf,
		history: history,
	}
}

func (c *chat) ChatRoomStream(ctx *gin.Context) {
	var request openai.ChatCompletionMessage
	err := ctx.ShouldBind(&request)
	if err != nil {
		Logger.Error(err)
		Response(ctx, http.StatusInternalServerError, "", err)
		return
	}

	content := c.history.GetValue(ctx.RemoteIP())
	content = append(content, request)
	resContent, err := TurboStream(ctx, &c.cnf.ChatGPT, content)
	if len(resContent) == 0 || err != nil {
		Response(ctx, 500, nil, err)
		return
	}

	newMsg := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleAssistant,
		Content: resContent,
	}
	c.history.SetValue(ctx.RemoteIP(), request)
	c.history.SetValue(ctx.RemoteIP(), newMsg)

	content = c.history.GetValue(ctx.RemoteIP())
	Logger.Infof("%+v", content)


	Response(ctx, 200, "", nil)
}


func (c *chat) ChatRoom(ctx *gin.Context) {
	var request openai.ChatCompletionMessage
	err := ctx.ShouldBind(&request)
	if err != nil {
		Logger.Error(err)
		Response(ctx, http.StatusInternalServerError, "", err)
		return
	}

	content := c.history.GetValue(ctx.RemoteIP())
	content = append(content, request)
	resContent, err := Turbo(ctx, &c.cnf.ChatGPT, content)

	Logger.Info(resContent)
	if len(resContent.Content) == 0 || err != nil {
		return
	}

	c.history.SetValue(ctx.RemoteIP(), request)
	c.history.SetValue(ctx.RemoteIP(), resContent)

	content = c.history.GetValue(ctx.RemoteIP())
	Logger.Infof("%+v", content)

	Response(ctx, 200, resContent.Content, nil)
}

func Response(ctx *gin.Context, code int, data interface{}, err error)  {
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": "",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  "success",
		"data": data,
	})
}