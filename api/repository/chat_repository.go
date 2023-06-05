package repository

import (
	"context"
	model "online_fashion_shop/api/model/chatbox"
	"online_fashion_shop/initializers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRepotitory interface {
	CreateDialog(context.Context, *model.Dialog) (*model.Dialog, error)
	CreateMessage(context.Context, *model.Message) (*model.Message, error)

	GetDialogByUserId(context.Context, string) (*model.Dialog, error)
	GetMessagesByDialogId(context.Context, string) ([]*model.Message, error)
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

func (r *chatRepotitory) GetMessagesByDialogId(ctx context.Context, dialogId string) (result []*model.Message, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	cursor, err := r.chatCollection.Find(ctx, bson.M{"dialog_id": dialogId})

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *chatRepotitory) GetDialogByUserId(ctx context.Context, userId string) (result *model.Dialog, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	rs := r.dialogCollection.FindOne(ctx, bson.M{"user_id": userId})

	err = rs.Decode(&result)

	return result, nil
}
