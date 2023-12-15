package main

import (
	"fmt"
	"net/http"

	"github.com/byteweap/weapon"
)

func main() {

	wp := weapon.New() // default config weapon
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wp.HandleRequest(w, r)
	})
	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		fmt.Println("fail ..... ", err.Error())
		return
	}
}
