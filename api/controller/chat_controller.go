package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"online_fashion_shop/api/common/errs"
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

func (c *ChatController) HandleWS(ctx *gin.Context) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		currentUser := ctx.MustGet("currentUser").(model.User)
		existed, err := c.Service.GetDialogByUserId(ctx, currentUser.Id)

		if err != nil {
			errs.HandleErrorStatus(ctx, err, "GetDialogByUserId")
			return
		}

		if existed == nil {
			existed, err = c.Service.CreateDialog(ctx)

			if err != nil {
				errs.HandleErrorStatus(ctx, err, "GetDialogByUserId")
				return
			}
		}
		c.s.Conns[ws] = currentUser.Id
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

// Chat
//
//	@Summary		Send user message
//	@Description    Send message
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Param          chat   			body        chat.CreateMessageInput    	true    "Message"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/send-message [post]
func (c *ChatController) SendMessage(ctx *gin.Context) {
	var data chat.CreateMessageInput

	if err := ctx.ShouldBind(&data); err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateMessage")
		return
	}

	rs, err := c.Service.CreateMessage(ctx, data)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateMessage")
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   *rs,
	})
}

// Chat
//
//	@Summary		Get user's messages
//	@Description    Get user's messages
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Param			page			query		int		true	"page"
//	@Param			page_size		query		int	 	true	"page size"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/message [get]
func (c *ChatController) GetUserMessage(ctx *gin.Context) {
	rs, err := c.Service.GetUserMessage(ctx)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserMessage")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   *rs,
	})
}

func (c *ChatController) SendUserMessage(ctx *gin.Context) {
	message := ctx.Param("message")
	dialogId := ctx.Param("dialogId")
	userId := ctx.Param("userId")

	rs, err := c.Service.CreateMessage(ctx, chat.CreateMessageInput{
		DialogId: dialogId,
		IsUser:   false,
		Type:     "text",
		Value:    message,
	})

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateMessage")
		return
	}

	data, err := json.Marshal(rs)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateMessage")
		return
	}

	for key, value := range c.s.Conns {
		if value == userId {
			go func(ws *websocket.Conn) {

				if _, err := ws.Write([]byte(data)); err != nil {
					fmt.Println("Write error: ", err)
				}
			}(key)

		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Send message successfully!",
	})
}

// Chat
//
//	@Summary		Get all dialog
//	@Description    Get all dialog
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/dialogs [get]
func (c *ChatController) GetAllDialogs(ctx *gin.Context) {
	rs, err := c.Service.GetDialogs(ctx)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetAllDialogs")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   rs,
	})
}
