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
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(rw, index)
	})); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown complete.")
}
