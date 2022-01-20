package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const index = `<html>
<h1>Joe's Music Project</h1>
</html>`

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("Starting up...")
	addr := ":" + os.Getenv("PORT")
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintln(rw, index)
		log.Println(r.URL, n, err)
	})); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown complete.")
}
