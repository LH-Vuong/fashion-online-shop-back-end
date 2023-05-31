package service

import (
	chat "online_fashion_shop/api/model/chatbox"
	model "online_fashion_shop/api/model/user"
	"time"

	"online_fashion_shop/api/repository"

	"github.com/gin-gonic/gin"
)

type ChatService interface {
	CreateDialog(*gin.Context, string) (*chat.Dialog, error)
}

type ChatServiceImpl struct {
	r repository.ChatRepotitory
}

func NewChatServiceImpl(r repository.ChatRepotitory) ChatService {
	return ChatServiceImpl{
		r: r,
	}
}

func (s ChatServiceImpl) CreateDialog(ctx *gin.Context, userId string) (*chat.Dialog, error) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	result, err := s.r.CreateDialog(ctx, &chat.Dialog{
		UserId:    currentUser.Id,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s ChatServiceImpl) CreateMessage(ctx *gin.Context, input chat.CreateMessageInput) (*chat.Message, error) {
	result, err := s.r.CreateMessage(ctx, &chat.Message{
		Value:     input.Value,
		DialogId:  input.DialogId,
		IsUser:    input.IsUser,
		CreatedAt: time.Now(),
		Type:      input.Type,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
