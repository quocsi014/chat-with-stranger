package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/quocsi014/util"
	"golang.org/x/net/websocket"
)


type User struct {
	name string
	conn *websocket.Conn
}

type Server struct {
	pair       map[*User]*User
	mu         sync.Mutex
	user_queue util.Queue
	disconnectedUser map[*User]bool
}

type Message struct {
	UserName string `json:"user_name"`
	IsSystem bool   `json:"is_system"`
	Message  string `json:"message"`
}

func NewServer() *Server {
	return &Server{
		pair:       make(map[*User]*User),
		mu:         sync.Mutex{},
		user_queue: *util.NewQueue(),
		disconnectedUser: make(map[*User]bool),
	}
}

func NewUser(name string, conn *websocket.Conn) *User {
	return &User{
		name: name,
		conn: conn,
	}
}

func NewMessage(userName, message string, isSystem bool) Message {
	return Message{
		UserName: userName,
		Message:  message,
		IsSystem: isSystem,
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {

	query := ws.Request().URL.Query()
	name := query.Get("name")

	if name == "" {
		name = "anonymous"
	}

	user := NewUser(name, ws)

	// iUser := s.user_queue.DeQueue()

	// if iUser == nil {
	// 	s.user_queue.EnQueue(user)

	// 	mes := NewMessage("System", "Pls, Wait another user", true)
	// 	if mesBytes, err := json.Marshal(mes); err != nil {
	// 		log.Fatal("Error encoding")
	// 	} else {
	// 		ws.Write(mesBytes)
	// 	}
	// } else {
	// 	waitingUser := iUser.(*User)
	// 	_, ok := s.disconnectedUser[waitingUser]
	// 	if ok{

	// 	}
	// 	s.pair[user] = waitingUser
	// 	s.pair[waitingUser] = user

	// 	mes := NewMessage("System", fmt.Sprintf("%s joined", user.name), true)
	// 	if mesBytes, err := json.Marshal(mes); err != nil {
	// 		log.Fatal("Error encoding")
	// 	} else {
	// 		waitingUser.conn.Write(mesBytes)
	// 	}
	// }
	s.CreateUserConn(user)
	s.readLoop(user)

}

func (s *Server)CreateUserConn(user *User){
	iUser := s.user_queue.DeQueue()

	if iUser == nil {
		s.user_queue.EnQueue(user)

		mes := NewMessage("System", "Pls, Wait another user", true)
		if mesBytes, err := json.Marshal(mes); err != nil {
			log.Fatal("Error encoding")
		} else {
			user.conn.Write(mesBytes)
		}
	} else {
		waitingUser := iUser.(*User)
		_, ok := s.disconnectedUser[waitingUser]
		if ok{
			s.CreateUserConn(user)
		}
		s.pair[user] = waitingUser
		s.pair[waitingUser] = user

		mes := NewMessage("System", fmt.Sprintf("%s joined", user.name), true)
		if mesBytes, err := json.Marshal(mes); err != nil {
			log.Fatal("Error encoding")
		} else {
			waitingUser.conn.Write(mesBytes)
		}
	}
}

func (s *Server) readLoop(user *User) {
	buf := make([]byte, 1024)
	for {
		n, err := user.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}

		mes := NewMessage(user.name, string(buf[:n]), false)
		if mesBytes, err := json.Marshal(mes); err != nil {
			log.Fatal("Error encoding")
		} else {
			receiveUser := s.pair[user]
			receiveUser.conn.Write(mesBytes)
		}

	}

	receiveUser := s.pair[user]
	if receiveUser == nil{
		s.disconnectedUser[user] = true
		return
	}


	mes := NewMessage("System", fmt.Sprintf("%s has left", user.name), true)
	if mesBytes, err := json.Marshal(mes); err != nil {
		log.Fatal("Error encoding")
	} else {
		receiveUser.conn.Write(mesBytes)
	}

	s.user_queue.EnQueue(s.pair[user])
	delete(s.pair, s.pair[user])
	delete(s.pair, user)

}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8080", nil)
}
