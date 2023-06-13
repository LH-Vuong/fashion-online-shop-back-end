package service

import (
	"online_fashion_shop/api/common/errs"
	chat "online_fashion_shop/api/model/chatbox"
	model "online_fashion_shop/api/model/user"
	"strconv"
	"time"

	"online_fashion_shop/api/repository"

	"github.com/gin-gonic/gin"
)

type ChatService interface {
	CreateDialog(*gin.Context) (*chat.Dialog, error)
	GetDialogByUserId(*gin.Context, string) (*chat.Dialog, error)
	CreateMessage(*gin.Context, chat.CreateMessageInput) (*chat.Message, error)
	GetUserMessage(*gin.Context) (*chat.GetUserMessageResponse, error)
}

type ChatServiceImpl struct {
	r repository.ChatRepotitory
}

func NewChatServiceImpl(r repository.ChatRepotitory) ChatService {
	return ChatServiceImpl{
		r: r,
	}
}

func (s ChatServiceImpl) CreateDialog(ctx *gin.Context) (*chat.Dialog, error) {
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

func (s ChatServiceImpl) GetDialogByUserId(ctx *gin.Context, userId string) (*chat.Dialog, error) {
	return s.r.GetDialogByUserId(ctx, userId)
}

func (s ChatServiceImpl) GetUserMessage(ctx *gin.Context) (*chat.GetUserMessageResponse, error) {
	currentUser := ctx.MustGet("currentUser").(model.User)
	pageSize, err := strconv.ParseInt(ctx.Query("page_size"), 10, 64)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserMessages")
		return nil, err
	}

	page, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserMessages")
		return nil, err
	}

	userDialog, err := s.r.GetDialogByUserId(ctx, currentUser.Id)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserMessages")
		return nil, err
	}

	userMessages, total, err := s.r.GetMessagesByDialogId(ctx, page, pageSize, userDialog.Id)

	if err != nil {
		return nil, err
	}
	return &chat.GetUserMessageResponse{
		DialogId: userDialog.Id,
		Messages: userMessages,
		Total:    total,
	}, nil
}
