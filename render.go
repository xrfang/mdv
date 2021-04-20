package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	gmt "github.com/mdigger/goldmark-toc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func RenderMD(fn string) (res map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = trace("%v", e)
		}
	}()
	f, err := os.Open(fn)
	assert(err)
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	assert(err)
	render := gmt.New(
		goldmark.WithExtensions(
			extension.DefinitionList,
			extension.Footnote,
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	var buf bytes.Buffer
	toc, err := render(src, &buf)
	assert(err)
	return map[string]interface{}{
		"toc": toc,
		"doc": strings.TrimSpace(buf.String()),
	}, nil
}
