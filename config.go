package weapon

import (
	"fmt"
	"time"
)

// Config weapon configuration.
type Config struct {
	WriteWait      time.Duration // Milliseconds until write times out.
	PongWait       time.Duration // Timeout for waiting on pong.
	PingPeriod     time.Duration // Milliseconds between pings.
	MaxMessageSize int64         // Maximum size in bytes of a message.

	connectHandler    func(*Session)
	closeHandler      func(*Session, int, string) error
	errorHandler      func(*Session, error)
	messageHandler    func(*Session, int, []byte)
	disconnectHandler func(*Session)
}

// 默认配置
func defaultConfig() *Config {
	return &Config{
		WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     time.Second * 9, //(60 * time.Second * 9) / 10,
		MaxMessageSize: 512,
		connectHandler: func(session *Session) {
			fmt.Println("connectHandler : ", session.ID(), " ", session.Request)
		},
		closeHandler: func(session *Session, i int, s string) error {
			fmt.Printf("closeHandler: id: %v, i: %v, s: %v \n", session.ID(), i, s)
			return nil
		},
		errorHandler: func(session *Session, err error) {
			if err != nil {
				fmt.Println("errorHandler: ", err.Error())
			}
		},
		messageHandler: func(session *Session, t int, bytes []byte) {
			fmt.Printf("Text ----- Id: %v Got msg: %v\n", session.ID(), string(bytes))
		},
		disconnectHandler: func(session *Session) {
			fmt.Println("disconnectHandler : ", session.ID())
		},
	}
}
