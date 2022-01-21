package main

import (
	"log"
	"net/http"
	"os"
)

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
