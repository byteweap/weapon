package weapon

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	id      int64
	Request *http.Request
	cfg     *Config
	conn    *websocket.Conn
	mux     *sync.RWMutex
	isOpen  bool
}

func newSession(id int64, cfg *Config, req *http.Request, conn *websocket.Conn) *Session {

	s := &Session{
		id:      id,
		Request: req,
		cfg:     cfg,
		conn:    conn,
		mux:     &sync.RWMutex{},
		isOpen:  true,
	}

	s.conn.SetReadLimit(s.cfg.MaxMessageSize)

	s.conn.SetReadDeadline(time.Now().Add(s.cfg.PongWait))
	s.conn.SetPongHandler(func(str string) error {
		s.conn.SetReadDeadline(time.Now().Add(s.cfg.PongWait))
		fmt.Println("pong : ", str)
		return nil
	})
	if s.cfg.closeHandler != nil {
		s.conn.SetCloseHandler(func(code int, text string) error {
			return s.cfg.closeHandler(s, code, text)
		})
	}

	return s
}

func (s *Session) ID() int64 {
	return s.id
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
}
