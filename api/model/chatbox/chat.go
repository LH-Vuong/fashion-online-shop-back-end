package chat

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	Conns map[string]*websocket.Conn
}

func (s *Server) broadcast(b []byte) {
	for key := range s.Conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Write error: ", err)
			}
		}(s.Conns[key])
	}
}

func NewServer() *Server {
	return &Server{
		Conns: make(map[string]*websocket.Conn),
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
}
