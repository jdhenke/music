package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("Starting up...")
	addr := ":" + os.Getenv("PORT")
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, http.HandlerFunc(serve)); err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown complete.")
}

var webFS = os.DirFS("web")

func serve(rw http.ResponseWriter, r *http.Request) {
	upath := filepath.Clean(r.URL.Path)
	if upath == "/" {
		upath = "."
	}
	upath = strings.TrimPrefix(upath, "/")
	f, err := webFS.Open(upath)
	if err != nil {
		handleErr(rw, err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	info, err := f.Stat()
	if err != nil {
		handleErr(rw, err)
		return
	}
	if info.IsDir() {
		upath = filepath.Join(upath, "index.html")
		i, err := webFS.Open(upath)
		if err != nil {
			handleErr(rw, err)
			return
		}
		_ = i.Close()
	}
	t := template.Must(template.New(filepath.Base(upath)).ParseFS(webFS, "templates/*", upath))
	if err := t.Execute(rw, nil); err != nil {
		handleErr(rw, err)
		return
	}
}

func handleErr(rw http.ResponseWriter, err error) {
	log.Println(err)
	if perr, ok := err.(*fs.PathError); ok && os.IsNotExist(perr.Err) {
		http.Error(rw, http.StatusText(404), 404)
		return
	}
	http.Error(rw, http.StatusText(500), 500)
}
