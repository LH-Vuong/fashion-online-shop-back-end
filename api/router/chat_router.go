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

		s.GET("api/open/:token", middleware.DeserializeUser(), controller.HandleWS)
		s.GET("api/message", middleware.DeserializeUser(), controller.GetUserMessage)
		s.POST("api/send-message", middleware.DeserializeUser(), controller.SendMessage)
		s.POST("api/send-user-message/:userId/:dialodId/:message", controller.SendUserMessage)
	})
	if err != nil {
		panic(err)
	}
}
