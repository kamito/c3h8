package propane

import (
	"os"
	"io/ioutil"
	"github.com/russross/blackfriday"
)


type Markdown struct {
	Source  []byte
	Content string
}

func (self Markdown) Render(path string) string {
	input := self.ReadFile(path)
	output := blackfriday.MarkdownBasic(input)
	return string(output)
}

func (self Markdown) ReadFile(path string) []byte {
	fp, _ := os.Open(path)
	defer fp.Close()
	buf, _ := ioutil.ReadAll(fp)
	return buf
}
