package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

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

type Server struct {
	rooms map[string]*Room
}

type Room struct {
	conns          map[*websocket.Conn]string
	anonymous_user int
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*Room),
	}
}

func NewRoom() *Room {
	return &Room{
		conns:          make(map[*websocket.Conn]string),
		anonymous_user: 0,
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {

	query := ws.Request().URL.Query()
	key := query.Get("room_key")
	name := query.Get("name")

	if key == "" {
		var errCreateKey error
		key, errCreateKey = s.createRoomKey()
		if errCreateKey != nil {
			ws.Write([]byte("Fail to create key, pls connect again or connect again with room key"))
		}
	}

	room, exist := s.rooms[key]

	if exist {
		if name == "" {
			name = fmt.Sprintf("anonymous%d", s.rooms[key].anonymous_user)
		}
		room.conns[ws] = name
		for key := range room.conns {
			key.Write([]byte(fmt.Sprintf("%s joined", name)))
		}
	} else {
		newRoom := NewRoom()
		if name == "" {
			name = "anonymous"
			newRoom.anonymous_user += 1

		}
		newRoom.conns[ws] = name
		s.rooms[key] = newRoom
		ws.Write([]byte(fmt.Sprintf("New room was created with key is %s", key)))
	}

	s.readLoop(ws, key)

}

func (s *Server) createRoomKey() (string, error) {
	for {
		key, err := generateRandomString(10)
		if err != nil {
			return "", err
		}

		_, exist := s.rooms[key]
		if !exist {
			return key, nil
		}

	}
}

func (s *Server) readLoop(ws *websocket.Conn, roomKey string) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}

		msg := buf[:n]
		s.broastcast(roomKey, msg, ws)
	}
}

func (s *Server)broastcast(roomKey string, msg []byte, sender *websocket.Conn){
	strMsg := fmt.Sprintf("%s: ", s.rooms[roomKey].conns[sender]) + string(msg)
	for ws := range s.rooms[roomKey].conns{
		ws.Write([]byte(strMsg))
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":3000", nil)
}
