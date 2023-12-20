package weapon

import (
	"fmt"
	"net/http"

	"github.com/byteweap/weapon/pkg/mapx"
	"github.com/gorilla/websocket"
)

type Weapon struct {
	config   *Config
	upgrader *websocket.Upgrader
	sessions mapx.Mapx[string, *Session] // 所有会话, key: 连接id
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
		sessions: mapx.New[string, *Session](true),
	}
}

// OneConnect 建立连接
func (w *Weapon) OneConnect(fn func(*Session)) {
	w.config.connectHandler = fn
}

// OnMessage 监听消息.
func (w *Weapon) OnMessage(fn func(*Session, int, []byte)) {
	w.config.messageHandler = fn
}

// OnDisconnect 断开链接.
func (w *Weapon) OnDisconnect(fn func(*Session)) {
	w.config.disconnectHandler = fn
}

// IdGenerator 自定义SessionId生成方法
// 默认以http.Request的Form表单中uid参数为sessionId,具体实现见defaultSessionIdGenerator.
func (w *Weapon) IdGenerator(fn func(*http.Request) string) {
	w.config.sessionIdGenerator = fn
}

func (w *Weapon) Run(pattern, addr string) {
	http.HandleFunc(pattern, func(responseWriter http.ResponseWriter, request *http.Request) {
		w.HandleRequest(responseWriter, request)
	})
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func (w *Weapon) Broadcast(msg []byte, filter func(s *Session) bool) {
	w.sessions.Range(func(_ string, s *Session) {
		if filter != nil && !filter(s) {
			return
		}
		s.push(msg)
	})
}

// 结合http使用
func (w *Weapon) HandleRequest(rw http.ResponseWriter, req *http.Request) {

	conn, err := w.upgrader.Upgrade(rw, req, rw.Header())
	if err != nil {
		fmt.Printf("Upgrade err: %v \n", err.Error())
		return
	}
	defer conn.Close()

	session, err := newSession(w.config, req, conn)
	if err != nil {
		fmt.Printf("NewSession err: %v \n", err.Error())
		return
	}
	defer session.close()

	w.config.connectHandler(session)
	session.ReadMessage()
	w.config.disconnectHandler(session)

	w.sessions.Delete(session.ID())
}
