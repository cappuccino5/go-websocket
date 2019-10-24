//
// Tips:
// 1. 客户端连接上后，需要马上上报自己的 id
// 2. 客户端在(0~5)分钟内至少需要ping一次服务器，否则服务器会主动断开连接, 开始回收connection
//
package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	defaultUpgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	defaultBinder = &Binder{
		userId2ConnMap: make(map[uint64]*Conn),
	}
)

type Server struct {
	upgrader       *websocket.Upgrader
	binder         *Binder
	RegisterHandle func(userId uint64)
}

func NewServer() *Server {
	server := &Server{upgrader: defaultUpgrader, binder: defaultBinder}
	return server
}

func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn := newConn(wsConn)
	var userId uint64
	conn.receiveHandle = func(data []byte, msgType int) {
		switch msgType {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			msg := Message{}
			err = json.Unmarshal(data, &msg)
			if err != nil || msg.From == 0 {
				return
			}
			if msg.Type == Type_Register {
				userId = msg.From
				oldConn, err := s.binder.ConnByUid(userId)
				if err == nil && oldConn != conn {
					oldConn.close()
				}
				s.binder.Bind(userId, conn)
				s.RegisterHandle(userId)
				return
			}
		}
	}
	conn.closeHandle = func() {
		if userId > 0 {
			s.binder.UnBind(userId)
		}
	}
}

func (s *Server) SendMessages(msgs []Message) error {
	if len(msgs) <= 0 {
		return nil
	}
	to := msgs[0].To
	if to == 0 {
		return ERR02
	}

	conn, err := s.binder.ConnByUid(to)
	if err != nil {
		return ERR03
	}

	data, err := json.Marshal(msgs)
	if err != nil {
		return ERR04
	}
	err = conn.sendMessage(websocket.BinaryMessage, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GetConnAll() (map[uint64]*Conn) {
	return s.binder.userId2ConnMap
}
