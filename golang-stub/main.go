package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{ID: 1, Name: "test"})
	})

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	addr := ":" + *port
	fmt.Printf("Starting server on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
