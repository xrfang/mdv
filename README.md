# mdv - Markdown Viewer

**MDV** is a "pure" markdown viewer in the sense that 1) it is just a viewer (not an editor); 2) it is a standalone tool, not a browser plugin.  Comparing to other viewers (or editors), **MDV** provides a pleasant viewing expience and at the same time is customizable and versatile.

## USAGE

### Installation

1. clone this repo and run `make` in the working directory.  To build for Windows, run `make windows`.
1. copy `mdv` (or `mdv.exe`) to a directory listed in system `PATH`, and associate it with `.md` file.

### Command Line

```none
Markdown Viewer V21.0945715

USAGE: mdv [OPTIONS] <markdown-file>

OPTIONS:

  -live
    	refresh when markdown file changes
  -port int
    	HTTP port (auto if not specified)
  -serve
    	serve only (do not open in browser)
  -version
    	show version info
  -wait
    	do not automatically quit server
```

### Customize

**MDV** will generate its configuration under `<USER-CONFIG>/mdv` folder.  `USER-CONFIG` is `~/.config` on Linux and `%APPDATA%` on Windows.  the `config.json` file in that folder contains the following options:

```json
{
    "css": "default.css",
    "port": 0,
    "quit": 9
}
```

Options are:

* **css**: name of the stylesheet to use
* **port**: HTTP port, 0 means auto select
* **quit**: delay before quit the local server

Usually, there is no need to edit this config file, except that you may want to use a new stylesheet.  Alternatively, you can also edit the `default.css` file directly to tweak display effects.   Also, there is a `highlight.css` which is the stylesheet for [syntax highlighting](https://highlightjs.org/).

### Other Usages

Apart from a simple viewer, **MDV** can also work as a "markdown preview engine", with the following options:

* `-live`: live update when markdown file changes
* `-port`: use a specific port instead of randomly choose one
* `-serve`: do not open a browser window (let the editor handle it)
* `-wait`: do not quit the server automatically

## Markdown Features

* [GoldMark supported features](https://github.com/yuin/goldmark):
  * Github-Flavored Markdown (Tables, Strikethrough, Linkify and TaskList)
  * DefinitionList
  * Footnotes
  * Typographer
  * Unsafe (parse embeded HTML tags)
* Table of Contents in collapsible side panel
* MathJax support - nicely show math formular
* Sytax highlighting for common programming languages

## Building Blocks

This tool utilizes the following excellent open source software:

* [yuin/goldmark ](https://github.com/yuin/goldmark): markdown parser for Go.
* [mdigger/goldmark-toc](https://github.com/mdigger/goldmark-toc): table of contents generator.
* [MathJax](https://www.mathjax.org/): beautiful and accessible math in all browsers. 
* [highlight.js](https://highlightjs.org/): syntax highlighting for the Web.
