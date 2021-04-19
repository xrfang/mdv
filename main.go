package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed resources/*
var res embed.FS

func main() {
	ver := flag.Bool("version", false, "show version info")
	flag.Usage = func() {
		fmt.Printf("Markdown Viewer %s\n\n", verinfo())
		fmt.Printf("USAGE: %s [OPTIONS] <markdown-file>\n\n", filepath.Base(os.Args[0]))
		fmt.Printf("OPTIONS:\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *ver {
		fmt.Println(verinfo())
		return
	}
	if len(flag.Args()) == 0 {
		fmt.Println("ERROR: markdown file not provided")
		os.Exit(1)
	}
	fn := flag.Arg(0)
	dir := filepath.Dir(fn)
	root, _ := fs.Sub(res, "resources")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, trace("%v", e).Error(), http.StatusInternalServerError)
			}
		}()
		f, err := root.Open(filepath.Join(".", r.URL.Path))
		if err == nil {
			defer f.Close()
			_, err = io.Copy(w, f)
			assert(err)
			return
		}
		if !os.IsNotExist(err) {
			panic(err)
		}
		f, err = os.Open(filepath.Join(dir, r.URL.Path))
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, fmt.Sprintf("not found: %s", r.URL.Path), http.StatusNotFound)
				return
			}
			panic(err)
		}
		defer f.Close()
		_, err = io.Copy(w, f)
		assert(err)
	})
	http.HandleFunc("/render", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(RenderMD(fn)))
	})
	ln, err := net.Listen("tcp", ":0")
	assert(err)
	port := ln.Addr().(*net.TCPAddr).Port
	go open(fmt.Sprintf("http://127.0.0.1:%d/index.html", port))
	panic(http.Serve(ln, nil))
}
