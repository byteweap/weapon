package weapon

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// session对象池
// 保存和复用临时对象,减少内存分配,降低GC压力
var sessionPool = &sync.Pool{
	New: func() any {
		return new(Session)
	},
}

type Session struct {
	id      string
	request *http.Request
	cfg     *Config
	conn    *websocket.Conn
	isOpen  bool
}

func newSession(cfg *Config, req *http.Request, conn *websocket.Conn) (*Session, error) {

	sid := cfg.sessionIdGenerator(req)
	if sid == "" {
		return nil, errors.New("session id is empty")
	}
	s := sessionPool.Get().(*Session)
	s.id = sid
	s.request = req
	s.cfg = cfg
	s.conn = conn
	s.isOpen = true

	s.conn.SetReadLimit(s.cfg.MaxMessageSize)
	// s.conn.SetReadDeadline(time.Now().Add(s.cfg.PongWait))
	// s.conn.SetPongHandler(func(str string) error {
	// 	s.conn.SetReadDeadline(time.Now().Add(s.cfg.PongWait))
	// 	return nil
	// })
	if s.cfg.closeHandler != nil {
		s.conn.SetCloseHandler(func(code int, text string) error {
			return s.cfg.closeHandler(s, code, text)
		})
	}
	return s, nil
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) HttpRequest() *http.Request {
	return s.request
}

func (s *Session) ReadMessage() {

	s.ping()

	for {
		t, message, err := s.conn.ReadMessage()
		if err != nil {
			s.cfg.errorHandler(s, err)
			break
		}
		s.cfg.messageHandler(s, t, message)
	}
}

func (s *Session) push(msg []byte) error {
	return s.conn.WriteMessage(websocket.TextMessage, msg)
}

func (s *Session) ping() {

	go func() {
		for {
			if !s.isOpen {
				return
			}
			time.Sleep(s.cfg.PingPeriod)
			s.conn.SetWriteDeadline(time.Now().Add(s.cfg.WriteWait))
			err := s.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				fmt.Println("ping err :", err.Error())
				break
			}
		}
	}()
}

func (s *Session) close() {
	s.isOpen = false
	s.putToPool()
}

func (s *Session) putToPool() {
	if s == nil {
		return
	}
	s.id = ""
	s.cfg = nil
	s.conn = nil
	s.request = nil
	s.isOpen = false

	sessionPool.Put(s)
}
