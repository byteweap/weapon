package weapon

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	id      string
	request *http.Request
	cfg     *Config
	conn    *websocket.Conn
	isOpen  bool
}

func newSession(cfg *Config, req *http.Request, conn *websocket.Conn) (*Session, error) {

	s := &Session{
		request: req,
		cfg:     cfg,
		conn:    conn,
		isOpen:  true,
	}
	s.id = cfg.sessionIdGenerator(req)
	if s.id == "" {
		return nil, errors.New("session id is empty")
	}
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
}
