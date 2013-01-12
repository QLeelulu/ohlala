package lib

import (
    "github.com/russross/blackfriday"
)

var renderer blackfriday.Renderer
var extensions = 0

func init() {

    // render the data into HTML
    var htmlFlags = 0
    htmlFlags |= blackfriday.HTML_SKIP_HTML
    htmlFlags |= blackfriday.HTML_SKIP_STYLE
    htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
    // htmlFlags |= blackfriday.HTML_USE_XHTML
    // htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
    // htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
    // htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
    // htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
    // htmlFlags |= blackfriday.HTML_OMIT_CONTENTS
    // htmlFlags |= blackfriday.HTML_TOC

    title := ""
    css := ""
    renderer = blackfriday.HtmlRenderer(htmlFlags, title, css)

    // set up options
    extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
    extensions |= blackfriday.EXTENSION_TABLES
    extensions |= blackfriday.EXTENSION_FENCED_CODE
    extensions |= blackfriday.EXTENSION_AUTOLINK
    extensions |= blackfriday.EXTENSION_STRIKETHROUGH
    extensions |= blackfriday.EXTENSION_SPACE_HEADERS
    extensions |= blackfriday.EXTENSION_HARD_LINE_BREAK // 强制换行

    // var output []byte
    // output = blackfriday.Markdown(input, renderer, extensions)
}

func Markdown(input []byte) []byte {
    return blackfriday.Markdown(input, renderer, extensions)
}
