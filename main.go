package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"

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
	addr := "0.0.0.0:8080"
	dir := "."

	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	fs := &CustomWebDAVFileSystem{FileSystem: webdav.Dir(dir)}

	handler := &webdav.Handler{
		FileSystem: fs,
		LockSystem: webdav.NewMemLS(),
	}

	log.Printf("Serving %s on %s", dir, addr)

	log.Fatal(http.ListenAndServe(addr, handler))
}

