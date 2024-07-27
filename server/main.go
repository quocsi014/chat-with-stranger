package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/quocsi014/util"
	"golang.org/x/net/websocket"
)

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode to base64 to get alphanumeric characters
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

type User struct{
	name string
	conn *websocket.Conn
}

type Server struct {
	pair map[*User]*User
	mu sync.Mutex
	user_queue util.Queue
}


type Message struct {
	UserName string `json:"user_name"`
	IsSystem bool `json:"is_system"`
	Message string `json:"message"`
}


func NewServer() *Server {
	return &Server{
		pair: make(map[*User]*User),
		mu: sync.Mutex{},
		user_queue: *util.NewQueue(),
	}
}

func NewUser(name string, conn *websocket.Conn) *User{
	return &User{
		name: name,
		conn: conn,
	}
}

func NewMessage(userName, message string, isSystem bool) Message{
	return Message{
		UserName: userName,
		Message: message,
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
	
	iUser := s.user_queue.DeQueue()

	if iUser == nil{
		s.user_queue.EnQueue(user)

		mes := NewMessage("System", "Pls, Wait another user", true)
		if mesBytes, err := json.Marshal(mes); err != nil{
			log.Fatal("Error encoding")
		}else{
			ws.Write(mesBytes)
		}
	}else{
		waitingUser := iUser.(*User)
		s.pair[user] = waitingUser
		s.pair[waitingUser] = user

		mes := NewMessage("System", fmt.Sprintf("%s joined", user.name), true)
		if mesBytes, err := json.Marshal(mes); err != nil{
			log.Fatal("Error encoding")
		}else{
			waitingUser.conn.Write(mesBytes)
		}
	}

	
	s.readLoop(user)

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
		if mesBytes, err := json.Marshal(mes); err != nil{
			log.Fatal("Error encoding")
		}else{
			receiveUser := s.pair[user]
			receiveUser.conn.Write(mesBytes)
		}

	}

	s.mu.Lock()
	mes := NewMessage(user.name, fmt.Sprintf("%s has left", user.name), false)
		if mesBytes, err := json.Marshal(mes); err != nil{
			log.Fatal("Error encoding")
		}else{
			receiveUser := s.pair[user]
			receiveUser.conn.Write(mesBytes)
		}
	s.mu.Unlock()

}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8080", nil)
}
