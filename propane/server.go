package propane

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
)

type Any interface{}

func NoError() error {
	return nil
}

type Server struct {
	Host string
	Port int
}

type Memo struct {
	Path  string
	Title string
}

type TemplateParams struct {
	Path  string
	Title string
	Body  string
	Files []Memo
}

func (self Server) Run() {
	bindAddr := fmt.Sprint(self.Host) + ":" + fmt.Sprint(self.Port)
	fmt.Println(fmt.Sprint(" [INFO] Start server with ", bindAddr))
	http.HandleFunc("/", self.Handler)
	http.ListenAndServe(bindAddr, nil)
}

func (self Server) Handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path := url.Path
	if path == "/" {
		self.HandleIndex(r, w)
	} else {
		if self.ServedFile(path, w, r) != true {
			self.HandlePage(r, w)
		}
	}
}

func (self Server) ServedFile(urlPath string, w http.ResponseWriter, r *http.Request) bool {
	fullpath := filepath.Join(CurDir(), urlPath)
	isDir, _ := IsDirectory(fullpath)
	if isDir == true {
		http.Error(w, "File not found", 404)
		return true
	} else {
		fileInfo, _ := os.Stat(fullpath)
		if fileInfo == nil {
			http.Error(w, "File not found", 404)
			return true
		}

		fileName := (fileInfo).Name()
		matched, _ := path.Match("*.md", fileName)
		if matched != true {
			fmt.Printf(" [INFO] %s\n", urlPath)
			http.ServeFile(w, r, fullpath)
			return true
		}
	}
	return false
}

func (self Server) HandleIndex(r *http.Request, w http.ResponseWriter) {
	files := []Memo{}
	params := &TemplateParams{
		Files: self.GetFiles("/", files),
		Title: CurDir(),
	}
	self.Render(w, params, []string{"assets/templates/layouts.html", "assets/templates/index.html"})
}

func (self Server) HandlePage(r *http.Request, w http.ResponseWriter) {
	url := r.URL
	path := url.Path
	fullpath, err := self.CheckFile(path)
	v := url.Query()
	if err == nil {
		fmt.Printf(" [INFO] %s\n", path)
		remarkFlag := v.Get("remark")
		if remarkFlag != "" {
			self.HandlePageRemark(path, fullpath, r, w)
		} else {
			self.HandlePageMarkdown(path, fullpath, r, w)
		}
	} else {
		fmt.Errorf("[ERROR] File not found: %s\n", path)
		http.Error(w, "File not found", 404)
	}
}

func (self Server) HandlePageMarkdown(path string, fullpath string, r *http.Request, w http.ResponseWriter) {
	markdown := new(Markdown)
	output := markdown.Render(fullpath)
	params := &TemplateParams{
		Path: path,
		Body: output,
	}
	self.Render(w, params, []string{"assets/templates/layouts.html", "assets/templates/page.html"})
}

func (self Server) HandlePageRemark(path string, fullpath string, r *http.Request, w http.ResponseWriter) {
	md := new(Markdown)
	output := string(md.ReadFile(fullpath))
	params := &TemplateParams{
		Path: path,
		Body: output,
	}
	// Load user template
	userTemplatePath := filepath.Join(CurDir(), "remark.html")
	stat, _ := os.Stat(userTemplatePath)
	if stat != nil {
		t := template.New("").Funcs(self.Helpers())
		t, _ = t.ParseFiles(userTemplatePath)
		t.ExecuteTemplate(w, "remark", params)
	} else {
		self.Render(w, params, []string{"assets/templates/page_remark.html"})
	}
}

func (self Server) GetFiles(targetPath string, files []Memo) []Memo {
	dir := filepath.Join(CurDir(), targetPath)
	isDir, _ := IsDirectory(dir)
	if isDir == true {
		fileInfos, _ := ioutil.ReadDir(dir)
		for _, fileInfo := range fileInfos {
			fileName := (fileInfo).Name()
			newPath := filepath.Join(targetPath, fileName)
			if fileInfo.IsDir() == true {
				files = self.GetFiles(newPath, files)
			} else {
				matched, _ := path.Match("*.md", fileName)
				if matched == true {
					title := self.GetTitle(newPath)
					memo := Memo{Path: newPath, Title: title}
					files = append(files, memo)
				}
			}
		}
		return files
	} else {
		emptyMemo := Memo{Title: "Directory is empty", Path: targetPath}
		files := []Memo{emptyMemo}
		return files
	}
}

func (self Server) GetTitle(targetPath string) string {
	filePath := filepath.Join(CurDir(), targetPath)
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		return ""
	}
	return line
}

func (self Server) Helpers() template.FuncMap {
	return template.FuncMap{
		"htmlsafe": func(context string) template.HTML {
			return template.HTML(context)
		},
	}
}

func (self Server) Render(w http.ResponseWriter, params Any, templates []string) {
	t := template.New("").Funcs(self.Helpers())
	for _, path := range templates {
		tmpl, err := Asset(path)
		if err == nil {
			t, _ = t.Parse(string(tmpl))
		}
	}
	// tmpls := template.Must(t.ParseFiles(templates...))
	t.ExecuteTemplate(w, "layouts", params)
}

func (self Server) CheckFile(path string) (string, error) {
	err := NoError()
	fullpath := filepath.Join(CurDir(), path)
	stat, _ := os.Stat(fullpath)
	if stat == nil {
		err = errors.New("File not found")
	}
	return fullpath, err
}

func RunServer(c *cli.Context) *Server {
	s := new(Server)
	s.Host = c.String("bind-addr")
	s.Port = c.Int("port")
	s.Run()
	return s
}
