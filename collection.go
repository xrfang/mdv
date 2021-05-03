package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func collect(root string) (col *collection, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	root, err = filepath.Abs(root)
	assert(err)
	st, err := os.Stat(root)
	assert(err)
	col = new(collection)
	col.Files = []string{}
	if st.IsDir() {
		col.Path = root
		col.Index = 0
		filepath.Walk(root, func(p string, fi os.FileInfo, e error) error {
			assert(err)
			if !fi.IsDir() && strings.ToLower(filepath.Ext(p)) == ".md" {
				col.Files = append(col.Files, p[len(col.Path)+1:])
			}
			return nil
		})
	} else {
		col.Path = filepath.Dir(root)
		col.Index = -1
		filepath.Walk(col.Path, func(p string, fi os.FileInfo, e error) error {
			assert(err)
			if !fi.IsDir() && strings.ToLower(filepath.Ext(p)) == ".md" {
				col.Files = append(col.Files, p[len(col.Path)+1:])
			}
			return nil
		})
	}
	sort.Slice(col.Files, func(i, j int) bool {
		return strings.ToLower(col.Files[i]) < strings.ToLower(col.Files[j])
	})
	if col.Index < 0 {
		bn := filepath.Base(root)
		for i, f := range col.Files {
			if f == bn {
				col.Index = i
				break
			}
		}
		if col.Index < 0 {
			col.Files = append([]string{bn}, col.Files...)
			col.Index = 0
		}
	}
	if len(col.Files) == 0 {
		col.Index = -1
	}
	return col, nil
}
