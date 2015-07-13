package propane

import (
	"io/ioutil"
	"os"
	// "github.com/russross/blackfriday"
	"github.com/shurcooL/github_flavored_markdown"
)

type Markdown struct {
	Source  []byte
	Content string
}

func (self Markdown) Render(path string) string {
	input := self.ReadFile(path)
	// output := blackfriday.MarkdownBasic(input)
	output := github_flavored_markdown.Markdown(input)
	return string(output)
}

func (self Markdown) ReadFile(path string) []byte {
	fp, _ := os.Open(path)
	defer fp.Close()
	buf, _ := ioutil.ReadAll(fp)
	return buf
}
