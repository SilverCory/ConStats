package web

import (
	"fmt"
	"net/http"
)

// RunWebserver Runs the server.
func RunWebserver(host string) {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/connectionData.json", dataHandler)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Println("An error occured whilst running the webserver..", err)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./connectionData.json")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}
