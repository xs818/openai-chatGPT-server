package router

import (
	"github.com/openai-chatGPT-server/config"
	"github.com/openai-chatGPT-server/logic"
	"github.com/openai-chatGPT-server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	logicHandler      logic.Chat
	config *config.Configuration
}

func NewRouter(config *config.Configuration, mh *middleware.Handler, logic logic.Chat) *Router {
	return &Router{
		config: config,
		logicHandler: logic,
	}
}

func (ru *Router) LoadRouters() http.Handler {

	gin.SetMode(ru.config.Server.Mode)

	r := gin.New()

	r.ForwardedByClientIP = true
	return ru.loadRouters(r)
}

func (ru *Router) loadRouters(r *gin.Engine) http.Handler {
	// Non-streaming Transfer
	r.POST("/chat-room", ru.logicHandler.ChatRoom)

	// Streaming Transfer
	r.POST("/chat-room/stream", ru.logicHandler.ChatRoomStream)
	return r
}