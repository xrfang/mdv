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
	"go.xrfang.cn/yal"
)

func RenderMD(fn string) (res map[string]interface{}, err error) {
	defer yal.Catch(&err)
	f, err := os.Open(fn)
	yal.Assert(err)
	defer f.Close()
	src, err := ioutil.ReadAll(f)
	yal.Assert(err)
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
	yal.Assert(err)
	return map[string]interface{}{
		"toc": toc,
		"doc": strings.TrimSpace(buf.String()),
	}, nil
}
