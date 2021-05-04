# mdv - Markdown Viewer

**MDV** is a "pure" markdown viewer in the sense that 1) it is just a viewer (not an editor); 2) it is a standalone tool, not a browser plugin. Comparing to other viewers (or editors), **MDV** provides a pleasant viewing experience and at the same time is customizable and versatile.

## USAGE

### Installation

1. Assuming you are using Linux or Mac, clone this repo and run `make` in the working directory. Windows users may download binary releases from [GitHub](https://github.com/xrfang/mdv/releases/) directly.
2. Cross-compiling is done by running `make windows`, `make linux` or `make mac` respectively.  
3. Copy `mdv` (or `mdv.exe`) to a directory listed in system `PATH`, and associate it with `.md` file.

### Command Line

```none
Markdown Viewer

USAGE: mdv [OPTIONS] <markdown-file>

OPTIONS:

  -port int
    	HTTP port (auto if not specified)
  -serve
    	serve only (do not open in browser)
  -version
    	show version info
```

### Customize

**MDV** will generate its configuration under `<USER-CONFIG>/mdv` folder. `USER-CONFIG` is `~/.config` on Linux, `%APPDATA%` on Windows and `~/Library/Application Support` on Mac.  The `config.json` file in that folder contains the following options:

```json
{
    "main_css": "default.css",
    "code_css": "highlight.css",
    "port":     0,
    "recurse":  2,
    "rev":      39
}
```

Options are:

* **main_css**: name of the primary stylesheet
* **code_css**: stylesheet for [syntax highlighting](https://highlightjs.org/)
* **port**: HTTP port, 0 means auto select
* **recurse**: level of sub-directories to search for markdown files (1~9, default 2)
* **rev**: revision of **MDV** that generated this config file

Usually, there is no need to edit this config file, except that you may want to use a new stylesheet. Alternatively, you can also edit the `default.css` file directly to tweak display effects. Also, there is a `highlight.css` which is the stylesheet for [syntax highlighting](https://highlightjs.org/).

> **Hint**: when the version of **MDV** upgrades, configuration will be reset.  If
> you use custom styles it is recommended that you use a new stylesheet and change
> `config.json` to point to your own style (rather than modify `default.css` directly).

### Other Usages

Apart from a simple viewer, **MDV** can also work as a "markdown preview engine", with the following options:

* `-port`: use a specific port instead of randomly choose one
* `-serve`: do not open a browser window (let the editor handle it)

## Markdown Features

* [GoldMark supported features](https://github.com/yuin/goldmark):
  * Github-Flavored Markdown (Tables, Strikethrough, Linkify and TaskList)
  * DefinitionList
  * Footnotes
  * Typographer
  * Unsafe (parse embeded HTML tags)
* Table of Contents in collapsible side panel
* File selector in collapsible side panel
* MathJax support - nicely show math formular
* Sytax highlighting for common programming languages

## Building Blocks

This tool utilizes the following excellent open source software:

* [yuin/goldmark ](https://github.com/yuin/goldmark): markdown parser for Go.
* [mdigger/goldmark-toc](https://github.com/mdigger/goldmark-toc): table of contents generator.
* [MathJax](https://www.mathjax.org/): beautiful and accessible math in all browsers. 
* [highlight.js](https://highlightjs.org/): syntax highlighting for the Web.
