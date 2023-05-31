package controller

import (
	"fmt"
	"io"
	chat "online_fashion_shop/api/model/chatbox"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type ChatController struct {
	Service service.ChatService
	s       chat.Server
}

func NewChatController(service service.ChatService, s chat.Server) *ChatController {
	return &ChatController{
		Service: service,
		s:       s,
	}
}

func (c *ChatController) CreateDialog(ctx *gin.Context) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		currentUser := ctx.MustGet("currentUser").(model.User)
		c.s.Conns[currentUser.Id] = ws
		buf := make([]byte, 1024)

		for {
			n, err := ws.Read(buf)

			if err != nil {
				if err == io.EOF {
					break
				}

				fmt.Println("read error:", err.Error())
			}

			msg := buf[:n]

			fmt.Println("recv:", string(msg))
		}

	})
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}

func (c *ChatController) SendMessage(ctx *gin.Context) {
	mess := ctx.Param("message")
	currentUser := ctx.MustGet("currentUser").(model.User)

	go func(ws *websocket.Conn) {
		if _, err := ws.Write([]byte(mess)); err != nil {
			fmt.Println("Write error: ", err)
		}
	}(c.s.Conns[currentUser.Id])
}
