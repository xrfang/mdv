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
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed resources/*
var res embed.FS

func main() {
	ver := flag.Bool("version", false, "show version info")
	serv := flag.Bool("serve", false, "serve only (do not open in browser)")
	port := flag.Int("port", 0, "HTTP port (auto if not specified)")
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
	if *port > 0 {
		cf.Port = *port
	}
	root, _ := fs.Sub(res, "resources")
	extract(root, "default.css")
	extract(root, "highlight.css")
	col, err := collect(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if col.Index < 0 {
		fmt.Fprintln(os.Stderr, "nothing to view")
		os.Exit(1)
	}
	var changed time.Time
	var mx sync.Mutex
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, trace("%v", e).Error(), http.StatusInternalServerError)
			}
		}()
		idx, err := strconv.Atoi(r.URL.Query().Get("idx"))
		if err == nil {
			col.Index = idx
		}
		entry := "/" + col.CurrentFile()
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, entry, http.StatusTemporaryRedirect)
			return
		}
		if r.URL.Path == entry {
			r.URL.Path = "/index.html"
		}
		switch strings.ToLower(filepath.Ext(r.URL.Path)) {
		case ".html", ".htm":
			w.Header().Set("Content-Type", "text/html")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
			switch r.URL.Path {
			case "/main.css":
				f, err := os.Open(filepath.Join(cf.dir, cf.MainCSS))
				assert(err)
				defer f.Close()
				_, err = io.Copy(w, f)
				assert(err)
				return
			case "/code.css":
				f, err := os.Open(filepath.Join(cf.dir, cf.CodeCSS))
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
		dir := filepath.Dir(col.CurrentPath())
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
		refreshTickCount()
		fn := col.CurrentPath()
		refresh := func() bool {
			mx.Lock()
			defer mx.Unlock()
			st, err := os.Stat(fn)
			assert(err)
			_, refresh := r.URL.Query()["refresh"]
			if refresh && !st.ModTime().After(changed) {
				return false
			}
			changed = st.ModTime()
			return true
		}()
		if !refresh {
			http.Error(w, "Not Modified", http.StatusNotModified)
			return
		}
		res, err := RenderMD(fn)
		assert(err)
		var fs []map[string]interface{}
		for i, f := range col.Files {
			dir := filepath.Dir(f)
			fn := filepath.Base(f)
			dir = strings.ReplaceAll(dir, string(os.PathSeparator), "/")
			if dir != "" {
				fn = path.Join(dir, fn)
			}
			fs = append(fs, map[string]interface{}{
				"idx": i,
				"sel": i == col.Index,
				"fn":  fn,
			})
		}
		res["col"] = fs
		w.Header().Set("Content-Type", "application/json")
		assert(json.NewEncoder(w).Encode(res))
	})
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cf.Port))
	assert(err)
	portInUse := ln.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d/", portInUse)
	fmt.Println("showing document at:", url)
	if !*serv {
		go open(url)
	}
	panic(http.Serve(ln, nil))
}
