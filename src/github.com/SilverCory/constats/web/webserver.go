package web

import (
	"fmt"
	"net/http"
)

// RunWebserver Runs the server.
func RunWebserver(host string) {
	err := http.ListenAndServe(host, http.FileServer(http.Dir("./data/")))
	if err != nil {
		fmt.Println("An error occured whilst running the webserver..", err)
	}
}
