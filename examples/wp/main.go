package main

import (
	"fmt"

	"github.com/byteweap/weapon"
)

func main() {

	wp := weapon.New()
	wp.OneConnect(func(s *weapon.Session) {
		fmt.Printf("New Conn: %v", s.ID())
	})
	wp.OnMessage(func(s *weapon.Session, i int, b []byte) {
		fmt.Printf("Type: %v, data: %v \n", i, string(b))
	})
	wp.OnDisconnect(func(s *weapon.Session) {
		fmt.Printf("Close Conn: %v", s.ID())
	})
	wp.Run("/ws", ":5001")
}
