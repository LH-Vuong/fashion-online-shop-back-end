package main

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"

// 	"golang.org/x/net/websocket"
// )

// type Server struct {
// 	conns map[*websocket.Conn]bool
// }

// func NewServer() *Server {
// 	return &Server{
// 		conns: make(map[*websocket.Conn]bool),
// 	}
// }

// func (s *Server) handleWS(ws *websocket.Conn) {
// 	fmt.Println("Add new connection: ", ws.RemoteAddr())

// 	s.conns[ws] = true

// 	s.readLoop(ws)
// }

// func (s *Server) handleWSNewChat(ws *websocket.Conn) {
// 	fmt.Println("New chat: ", ws.RemoteAddr())

// 	for {
// 		payload := fmt.Sprintf("New chat at %d\n ", time.Now().UnixNano())
// 		ws.Write([]byte(payload))
// 		time.Sleep(time.Second * 2)
// 	}
// }

// func (s *Server) readLoop(ws *websocket.Conn) {
// 	buf := make([]byte, 1024)

// 	for {
// 		n, err := ws.Read(buf)

// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}

// 			fmt.Println("read error:", err.Error())
// 		}

// 		msg := buf[:n]

// 		fmt.Println("recv:", string(msg))

// 		s.broadcast(msg)
// 	}
// }

// func main() {
// 	server := NewServer()

// 	http.Handle("/ws", websocket.Handler(server.handleWS))

// 	http.Handle("/new-chat", websocket.Handler(server.handleWSNewChat))

// 	http.ListenAndServe(":3000", nil)
// }
