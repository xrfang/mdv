package main

import (
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"go.xrfang.cn/yal"
)

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// extract embeded resource to config directory
func extract(root fs.FS, fn string) {
	fp := filepath.Join(cf.dir, fn)
	_, err := os.Stat(fp)
	if err != nil {
		func() {
			w, err := os.Create(fp)
			yal.Assert(err)
			defer func() { yal.Assert(w.Close()) }()
			f, _ := root.Open(filepath.Join(".", fn))
			defer f.Close()
			_, err = io.Copy(w, f)
			yal.Assert(err)
		}()
	}
}
