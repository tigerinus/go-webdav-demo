package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

type CustomWebDAVFileSystem struct {
	webdav.FileSystem
}

func (cfs *CustomWebDAVFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	info, err := cfs.FileSystem.Stat(ctx, name)
	if err != nil {
		return nil, err
	}
	return &CustomFileInfo{FileInfo: info}, nil
}

type CustomFileInfo struct {
	os.FileInfo
}

func (cfi *CustomFileInfo) Name() string {
	name := cfi.FileInfo.Name()
	return fmt.Sprintf("%s [First 3: %s]", name, firstThreeLetters(name))
}

func firstThreeLetters(s string) string {
	if len(s) < 3 {
		return s
	}
	return s[:3]
}

func main() {
	addr := "0.0.0.0:8888"

	var fs webdav.FileSystem

	if len(os.Args) < 2 {
		log.Printf("Serving in-memory filesystem on %s", addr)
		fs = webdav.NewMemFS()
	} else {
		dir := os.Args[1]
		log.Printf("Serving %s on %s", dir, addr)
		fs = &CustomWebDAVFileSystem{FileSystem: webdav.Dir(dir)}
	}

	handler := &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
	}

	log.Fatal(http.ListenAndServe(addr, handler))
}
