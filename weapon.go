package weapon

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Weapon struct {
	config   *Config
	upgrader *websocket.Upgrader
	sessions map[int64]*Session // 所有连接, key: 连接id
}

func New() *Weapon {

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return &Weapon{
		config:   defaultConfig(),
		upgrader: upgrader,
		sessions: make(map[int64]*Session),
	}
}

// OneConnect 建立连接
func (w *Weapon) OneConnect(fn func(*Session)) {
	w.config.connectHandler = fn
}

// OnMessage 监听消息
func (w *Weapon) OnMessage(fn func(*Session, int, []byte)) {
	w.config.messageHandler = fn
}

// OnDisconnect 断开链接
func (w *Weapon) OnDisconnect(fn func(*Session)) {
	w.config.disconnectHandler = fn
}

func (w *Weapon) Run(pattern, addr string) {
	http.HandleFunc(pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		w.HandleRequest(responseWriter, request)
	})
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

// 结合http使用
func (w *Weapon) HandleRequest(rw http.ResponseWriter, req *http.Request) error {
	conn, err := w.upgrader.Upgrade(rw, req, rw.Header())
	if err != nil {
		return err
	}
	session := newSession(666, w.config, req, conn) // todo Id如何生成
	w.config.connectHandler(session)

	session.ReadMessage()

	session.isOpen = false
	w.config.disconnectHandler(session)
	return nil
}
