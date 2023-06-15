package repository

import (
	"context"
	chat "online_fashion_shop/api/model/chatbox"
	model "online_fashion_shop/api/model/chatbox"
	"online_fashion_shop/initializers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepotitory interface {
	CreateDialog(context.Context, *model.Dialog) (*model.Dialog, error)
	CreateMessage(context.Context, *model.Message) (*model.Message, error)

	GetDialogByUserId(context.Context, string) (*model.Dialog, error)
	GetMessagesByDialogId(context.Context, int64, int64, string) ([]*model.Message, int64, error)
	GetAllDialogs(context.Context) ([]*chat.Dialog, error)
	GetDialogLatestMessage(context.Context, string) (*chat.Message, error)
}

type chatRepotitory struct {
	chatCollection   initializers.Collection
	dialogCollection initializers.Collection
}

func NewChatRepotitory(chatCollection initializers.Collection, dialogCollection initializers.Collection) ChatRepotitory {
	return &chatRepotitory{
		chatCollection:   chatCollection,
		dialogCollection: dialogCollection,
	}
}

func (r *chatRepotitory) CreateDialog(ctx context.Context, newDialog *model.Dialog) (*model.Dialog, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.dialogCollection.InsertOne(ctx, newDialog)

	if err != nil {
		return nil, err
	}

	newDialog.Id = res.(primitive.ObjectID).Hex()
	return newDialog, err
}

func (r *chatRepotitory) CreateMessage(ctx context.Context, newMessage *model.Message) (*model.Message, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.chatCollection.InsertOne(ctx, newMessage)

	if err != nil {
		return nil, err
	}

	newMessage.Id = res.(primitive.ObjectID).Hex()
	return newMessage, err
}

func (r *chatRepotitory) GetMessagesByDialogId(ctx context.Context, page int64, pageSize int64, dialogId string) (result []*model.Message, total int64, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	opts := options.Find().
		SetSkip(page - 1).
		SetLimit(pageSize)

	query := bson.M{"dialog_id": dialogId}

	rs, err := r.chatCollection.Find(ctx, query, opts)

	if err != nil {
		return nil, 0, err
	}

	total, err = r.chatCollection.CountDocuments(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	err = rs.All(ctx, &result)

	if err != nil {
		return nil, 0, err
	}
	return
}

func (r *chatRepotitory) GetDialogByUserId(ctx context.Context, userId string) (result *model.Dialog, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	rs := r.dialogCollection.FindOne(ctx, bson.M{"user_id": userId})

	err = rs.Decode(&result)

	return result, nil
}

func (r *chatRepotitory) GetAllDialogs(ctx context.Context) ([]*chat.Dialog, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	var result []*chat.Dialog

	rs, err := r.dialogCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	if err := rs.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *chatRepotitory) GetDialogLatestMessage(ctx context.Context, dialogId string) (*chat.Message, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	var result chat.Message

	options := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	if err := r.chatCollection.FindOne(ctx, bson.M{"dialog_id": dialogId}, options).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
