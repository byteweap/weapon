package main

import (
	"fmt"
	"net/http"

	"github.com/byteweap/weapon"
)

func main() {

	wp := weapon.New()
	wp.OneConnect(func(s *weapon.Session) {
		fmt.Printf("New Conn: %v \n", s.ID())
	})
	wp.OnMessage(func(_ *weapon.Session, i int, b []byte) {
		fmt.Printf("Type: %v, data: %v \n", i, string(b))
	})
	wp.OnDisconnect(func(s *weapon.Session) {
		fmt.Printf("Close Conn: %v \n", s.ID())
	})
	wp.IdGenerator(func(r *http.Request) string {
		return r.FormValue("uid1")
	})
	wp.Run("/ws", ":5001")
}
