package propane

import (
	"os"
	"fmt"
	"errors"
	"path/filepath"
	"net/http"
	"html/template"
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
}

type Server struct {
	Host string
	Port int
}


func (self Server) Run() {
	bindAddr := fmt.Sprint(self.Host) + ":" + fmt.Sprint(self.Port)
	http.HandleFunc("/", self.Handler)
	http.ListenAndServe(bindAddr, nil)
}

func (self Server) Handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path := url.Path
	if path == "/" {
		self.HandleIndex(w)
	} else {
		self.HandlePage(path, w)
	}
}

func (self Server) HandleIndex(w http.ResponseWriter) {
	params := &TemplateParams{
		Body: "Hello, C3H8 with index",
	}
	self.Render(w, params, []string{"assets/templates/index.html"})
}

func (self Server) HandlePage(path string, w http.ResponseWriter) {
	fullpath, err := self.CheckFile(path)
	if err == nil {
		fmt.Printf(" [INFO] %s\n", path)
		markdown := new(Markdown)
		output := markdown.Render(fullpath)
		params := &TemplateParams{
			Path: path,
			Body: output,
		}
		self.Render(w, params, []string{"assets/templates/page.html"})
	} else {
		fmt.Errorf("[ERROR] File not found: %s\n", path)
		http.Error(w, "File not found", 404)
	}
}

func (self Server) Helpers() template.FuncMap {
	return template.FuncMap{
		"htmlsafe": func(context string) template.HTML {
			return template.HTML(context)
		},
	}
}

func (self Server) Render(w http.ResponseWriter, params Any, templatePath []string) {
	templates := append([]string{"assets/templates/layouts.html"}, templatePath...)
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
	current, _ := filepath.Abs(".")
	fullpath := filepath.Join(current, path)
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
