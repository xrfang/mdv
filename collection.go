package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.xrfang.cn/yal"
)

type (
	collection struct {
		Path  string
		Files []string
		Index int
	}
)

func (c collection) CurrentFile() string {
	return c.Files[c.Index]
}

func (c collection) CurrentPath() string {
	return filepath.Join(c.Path, c.CurrentFile())
}

func iter(path string, level int) []string {
	d, err := os.Open(path)
	yal.Assert(err)
	defer d.Close()
	fis, err := d.ReadDir(0)
	yal.Assert(err)
	sub := make(map[string][]string)
	var items []string
	for _, fi := range fis {
		fp := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			if !strings.HasSuffix(fp, string(os.PathSeparator)) {
				fp += string(os.PathSeparator)
			}
			items = append(items, fp)
			if level > 0 {
				sub[fp] = iter(fp, level-1)
			}
			continue
		}
		ext := strings.ToLower(filepath.Ext(fi.Name()))
		if ext == ".md" {
			items = append(items, fp)
		}
	}
	sort.Strings(items)
	var files []string
	for _, fn := range items {
		if strings.HasSuffix(fn, string(os.PathSeparator)) {
			files = append(files, sub[fn]...)
			continue
		}
		files = append(files, fn)
	}
	return files
}

func collect(entry string, depth int) (col *collection, err error) {
	if depth < 1 || depth > 9 {
		depth = 2
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	entry, err = filepath.Abs(entry)
	yal.Assert(err)
	root := entry
	st, err := os.Stat(entry)
	yal.Assert(err)
	if !st.IsDir() {
		root = filepath.Dir(entry)
	}
	if !strings.HasSuffix(root, string(os.PathSeparator)) {
		root += string(os.PathSeparator)
	}
	files := iter(root, depth-1)
	var idx int
	if len(files) == 0 {
		idx = -1
	} else {
		for i := range files {
			if files[i] == entry {
				idx = i
			}
			files[i] = files[i][len(root):]
		}
	}
	return &collection{
		Path:  root,
		Files: files,
		Index: idx,
	}, nil
}
