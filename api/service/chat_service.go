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
	GetDialogs(*gin.Context) ([]*chat.GetDialogsResponse, error)
}

type ChatServiceImpl struct {
	chatRepo repository.ChatRepotitory
	userRepo repository.UserRepository
}

func NewChatServiceImpl(chatRepo repository.ChatRepotitory, userRepo repository.UserRepository) ChatService {
	return ChatServiceImpl{
		chatRepo: chatRepo,
		userRepo: userRepo,
	}
}

func (s ChatServiceImpl) CreateDialog(ctx *gin.Context) (*chat.Dialog, error) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	result, err := s.chatRepo.CreateDialog(ctx, &chat.Dialog{
		UserId:    currentUser.Id,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s ChatServiceImpl) CreateMessage(ctx *gin.Context, input chat.CreateMessageInput) (*chat.Message, error) {
	result, err := s.chatRepo.CreateMessage(ctx, &chat.Message{
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
	return s.chatRepo.GetDialogByUserId(ctx, userId)
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

	userDialog, err := s.chatRepo.GetDialogByUserId(ctx, currentUser.Id)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserMessages")
		return nil, err
	}

	userMessages, total, err := s.chatRepo.GetMessagesByDialogId(ctx, page, pageSize, userDialog.Id)

	if err != nil {
		return nil, err
	}
	return &chat.GetUserMessageResponse{
		DialogId: userDialog.Id,
		Messages: userMessages,
		Total:    total,
	}, nil
}

func (s ChatServiceImpl) GetDialogs(ctx *gin.Context) ([]*chat.GetDialogsResponse, error) {
	var result []*chat.GetDialogsResponse

	dialogList, err := s.chatRepo.GetAllDialogs(ctx)

	if err != nil {
		return nil, err
	}

	if dialogList == nil {
		return result, nil
	}

	for _, dialog := range dialogList {
		message, err := s.chatRepo.GetDialogLatestMessage(ctx, dialog.Id)

		if err != nil {
			return result, err
		}

		if message == nil {
			continue
		}

		user, err := s.userRepo.GetUserById(ctx, dialog.UserId)

		if err != nil {
			return result, err
		}

		if user == nil {
			continue
		}

		result = append(result, &chat.GetDialogsResponse{
			Id:           dialog.Id,
			UserId:       user.Id,
			UserFullname: user.Fullname,
			UserPhoto:    user.Photo,
			Message:      *message,
		})
	}

	return result, nil
}
