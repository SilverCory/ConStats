package web

import (
	"fmt"
	"net/http"
)

// RunWebserver Runs the server.
func RunWebserver(host string) {
	http.FileServer(http.Dir("./data/"))
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Println("An error occured whilst running the webserver..", err)
	}
}
