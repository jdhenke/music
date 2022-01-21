package main

import (
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
	if err := http.ListenAndServe(addr, http.FileServer(http.Dir("web"))); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown complete.")
}
