package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func RenderMD(fn string) (content string) {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			if os.IsNotExist(err) {
				content = ""
			} else {
				content = err.Error()
			}
		}
	}()
	f, err := os.Open(fn)
	assert(err)
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	assert(err)
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.DefinitionList,
			extension.Footnote,
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	var buf bytes.Buffer
	assert(md.Convert(src, &buf))
	return strings.TrimSpace(buf.String())
}
