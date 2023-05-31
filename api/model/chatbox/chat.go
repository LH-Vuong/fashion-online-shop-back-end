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
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Id        string    `json:"id"`
	DialogId  string    `json:"dialog_id"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	IsUser    bool      `json:"is_user"`
}

type CreateMessageInput struct {
	DialogId string `json:"dialog_id"`
	Value    string `json:"value"`
	IsUser   bool   `json:"is_user"`
	Type     string `json:"type"`
}
