package weapon

import (
	"fmt"
	"net/http"
	"time"
)

// Config weapon configuration.
type Config struct {
	WriteWait      time.Duration // Milliseconds until write times out.
	PongWait       time.Duration // Timeout for waiting on pong.
	PingPeriod     time.Duration // Milliseconds between pings.
	MaxMessageSize int64         // Maximum size in bytes of a message.

	connectHandler     func(*Session)
	closeHandler       func(*Session, int, string) error
	errorHandler       func(*Session, error)
	messageHandler     func(*Session, int, []byte)
	disconnectHandler  func(*Session)
	sessionIdGenerator func(*http.Request) string
}

var (
	defaultConnectHandler = func(session *Session) {
		fmt.Println("connectHandler : ", session.ID(), " ", session.HttpRequest())
	}
	defaultCloseHandler = func(session *Session, i int, s string) error {
		fmt.Printf("closeHandler: id: %v, i: %v, s: %v \n", session.ID(), i, s)
		return nil
	}
	defaultErrorHandler = func(_ *Session, err error) {
		if err != nil {
			fmt.Println("errorHandler: ", err.Error())
		}
	}
	defaultMessageHandler = func(session *Session, i int, bytes []byte) {
		fmt.Printf("Text ----- Id: %v Got msg: type: %v, data: %v\n", session.ID(), i, string(bytes))
	}
	defaultDisconnectionHandler = func(session *Session) {
		fmt.Println("disconnectHandler : ", session.ID())
	}
	defaultSessionIdGenerator = func(req *http.Request) string {
		return req.FormValue("uid")
	}
)

// 默认配置
func defaultConfig() *Config {
	return &Config{
		WriteWait:          10 * time.Second,
		PongWait:           10 * time.Second,
		PingPeriod:         time.Second * 9, //(60 * time.Second * 9) / 10,
		MaxMessageSize:     512,
		connectHandler:     defaultConnectHandler,
		closeHandler:       defaultCloseHandler,
		errorHandler:       defaultErrorHandler,
		messageHandler:     defaultMessageHandler,
		disconnectHandler:  defaultDisconnectionHandler,
		sessionIdGenerator: defaultSessionIdGenerator,
	}
}
