package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func main() {
	addr := "0.0.0.0:8888"

	var fs webdav.FileSystem

	if len(os.Args) < 2 {
		log.Printf("Serving in-memory filesystem on %s", addr)
		fs = webdav.NewMemFS()
	} else {
		dir := os.Args[1]
		log.Printf("Serving %s on %s", dir, addr)
		fs = webdav.Dir(dir)
	}

	handler := &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.Printf("%s %s: %v", r.Method, r.URL.Path, err)
		},
	}

	log.Print(http.ListenAndServe(addr, handler))
}
