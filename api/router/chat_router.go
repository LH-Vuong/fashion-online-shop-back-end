package router

import (
	"online_fashion_shop/api/controller"
	middleware "online_fashion_shop/api/middlewares"
	chat "online_fashion_shop/api/model/chatbox"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"golang.org/x/net/websocket"
)

func InitChatRouter(s *gin.Engine, c *dig.Container) {
	err := c.Invoke(func(chatService service.ChatService) {
		controller := controller.NewChatController(chatService, chat.Server{
			Conns: make(map[string]*websocket.Conn),
		})

		s.GET("api/open/:token", middleware.DeserializeUser(), controller.CreateDialog)

		s.GET("api/send/:message", middleware.DeserializeUser(), controller.SendMessage)
	})
	if err != nil {
		panic(err)
	}
}
