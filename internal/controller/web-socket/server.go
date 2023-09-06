package web_socket

import (
	"log"
	"net/http"
)

// Run starts a new chat server with 4 chat rooms, listening on port 8080
func Run() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	for _, name := range []string{"arduino", "java", "go", "scala"} {
		r := NewRoom(name)
		http.Handle("/chat/"+name, r)
		go r.Run()
	}

	log.Printf("starting chat server on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("starting server failed: ", err)
	}

}
