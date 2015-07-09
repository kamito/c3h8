package propane

import (
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

type TemplateParams struct {
	Path  string
	Title string
	Body  string
	Files []string
}

type Server struct {
	Host string
	Port int
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
		self.HandlePage(r, w)
	}
}

func (self Server) HandleIndex(r *http.Request, w http.ResponseWriter) {
	params := &TemplateParams{
		Files: self.GetFiles(),
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

func (self Server) GetFiles() []string {
	curDir := CurDir()
	isDir, _ := IsDirectory(curDir)
	if isDir == true {
		files := []string{}
		fileInfos, _ := ioutil.ReadDir(curDir)
		for _, fileInfo := range fileInfos {
			fileName := (fileInfo).Name()
			matched, _ := path.Match("*.md", fileName)
			if matched == true {
				files = append(files, fileName)
			}
		}
		return files
	} else {
		files := []string{"Directory is empty"}
		return files
	}
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
