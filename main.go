package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	root, _ := fs.Sub(res, "resources")
	defcss := filepath.Join(cf.dir, "default.css")
	_, err := os.Stat(defcss)
	if err != nil {
		func() {
			w, err := os.Create(defcss)
			assert(err)
			defer func() { assert(w.Close()) }()
			f, _ := root.Open(filepath.Join(".", "default.css"))
			defer f.Close()
			_, err = io.Copy(w, f)
			assert(err)
		}()
	}
	fn := flag.Arg(0)
	dir := filepath.Dir(fn)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, trace("%v", e).Error(), http.StatusInternalServerError)
			}
		}()
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect)
			return
		}
		switch strings.ToLower(filepath.Ext(r.URL.Path)) {
		case ".html", ".htm":
			w.Header().Set("Content-Type", "text/html")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
			if r.URL.Path == "/mdv.css" {
				f, err := os.Open(filepath.Join(cf.dir, cf.CSS))
				assert(err)
				defer f.Close()
				_, err = io.Copy(w, f)
				assert(err)
				return
			}
		case ".js":
			w.Header().Set("Content-Type", "text/javascript")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".png":
			w.Header().Set("Content-Type", "text/png")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}
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
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, trace("%v", e).Error(), http.StatusInternalServerError)
			}
		}()
		res, err := RenderMD(fn)
		assert(err)
		w.Header().Set("Content-Type", "application/json")
		assert(json.NewEncoder(w).Encode(res))
	})
	ln, err := net.Listen("tcp", ":0")
	assert(err)
	port := ln.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d/", port)
	fmt.Println("showing document at:", url)
	go func() {
		open(url)
		if cf.Quit > 0 {
			fmt.Printf("quit local server after %d seconds", cf.Quit)
			for i := 0; i < cf.Quit; i++ {
				time.Sleep(time.Second)
				fmt.Print(".")
			}
			fmt.Println(" bye")
			os.Exit(0)
		}
	}()
	panic(http.Serve(ln, nil))
}
