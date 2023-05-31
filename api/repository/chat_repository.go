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

	GetMessagesByUserId(context.Context, string) ([]*model.Message, error)
}

type chatRepotitory struct {
	chatCollection initializers.Collection
}

func NewChatRepotitory(chatCollection initializers.Collection) ChatRepotitory {
	return &chatRepotitory{
		chatCollection: chatCollection,
	}
}

func (r *chatRepotitory) CreateDialog(ctx context.Context, newDialog *model.Dialog) (*model.Dialog, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.chatCollection.InsertOne(ctx, newDialog)

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

func (r *chatRepotitory) GetMessagesByUserId(ctx context.Context, userId string) (result []*model.Message, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	rs := r.chatCollection.FindOne(ctx, bson.M{"user_id": userId})

	err = rs.Decode(&result)

	return result, nil
}
