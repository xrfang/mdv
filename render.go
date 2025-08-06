package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	gmt "github.com/mdigger/goldmark-toc"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"go.xrfang.cn/act"
)

func RenderMD(fn string) (res map[string]interface{}, err error) {
	defer act.Catch(&err)
	src, err := os.ReadFile(fn)
	act.Assert(err)
	base := filepath.Base(fn)
	ext := filepath.Ext(base)
	base = base[:len(base)-len(ext)]
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
	act.Assert(err)
	return map[string]interface{}{
		"toc":   toc,
		"title": base,
		"doc":   strings.TrimSpace(buf.String()),
	}, nil
}
