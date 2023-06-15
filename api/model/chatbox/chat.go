package chat

import (
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	Conns map[*websocket.Conn]string
}

func NewServer() *Server {
	return &Server{
		Conns: make(map[*websocket.Conn]string),
	}
}

type Dialog struct {
	Id        string    `bson:"_id,omitempty" json:"id"`
	UserId    string    `bson:"user_id" json:"user_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Message struct {
	Id        string    `bson:"_id,omitempty" json:"id"`
	DialogId  string    `bson:"dialog_id" json:"dialog_id"`
	Value     string    `bson:"value" json:"value"`
	Type      string    `bson:"type" json:"type"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	IsUser    bool      `bson:"is_user" json:"is_user"`
}

type CreateMessageInput struct {
	DialogId string `json:"dialog_id"`
	Value    string `json:"value"`
	IsUser   bool   `json:"is_user"`
	Type     string `json:"type"`
}

type GetUserMessageResponse struct {
	DialogId string     `json:"dialog_id"`
	Messages []*Message `json:"messages"`
	Total    int64      `json:"total"`
}

type GetDialogsResponse struct {
	Id           string  `json:"id"`
	UserId       string  `json:"user_id"`
	UserFullname string  `json:"user_fullname"`
	UserPhoto    string  `json:"user_photo"`
	Message      Message `json:"message"`
}
