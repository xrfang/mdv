package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO WORLD"))
	})
	ln, err := net.Listen("tcp", ":0")
	assert(err)
	port := ln.Addr().(*net.TCPAddr).Port
	fmt.Println("Using port:", port)
	go open(fmt.Sprintf("http://127.0.0.1:%d/", port))
	panic(http.Serve(ln, nil))
}
