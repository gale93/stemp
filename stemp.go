package stemp

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Stemplate is the main object of the package.
// The users will interface only with this one.
type Stemplate struct {
	// LiveReload set to TRUE make the changes to html pages aviable without restarting the server
	// [Default]: false
	LiveReload bool

	templatesDir string
	templates    map[string]*template.Template
}

//NewStemplate create a new Istance of Stemplate obj
func NewStemplate(templatesDirectory string) (*Stemplate, error) {
	var st Stemplate

	st.templatesDir = templatesDirectory
	st.LiveReload = false
	st.templates = make(map[string]*template.Template)

	if err := st.load(); err != nil {
		return nil, err
	}

	return &st, nil
}

// load function will store in RAM the , just compiled, templates.
func (st *Stemplate) load() error {

	templates, terr := filepath.Glob(st.templatesDir + "*.tmpl")
	if terr != nil {
		return terr
	}

	contents, err := filepath.Glob(st.templatesDir + "*.html")
	if err != nil {
		return err
	}

	for _, c := range contents {
		current := append(templates, c)
		st.templates[filepath.Base(c)] = template.Must(template.ParseFiles(current...))
	}

	return nil

}

// Reload is an utily function that allow to recompile templates at run time
func (st *Stemplate) Reload() {

	stReloaded, err := NewStemplate(st.templatesDir)
	if err != nil {
		fmt.Println("[Reloading templates]: ERROR")
		return
	}

	stReloaded.LiveReload = st.LiveReload

	st = stReloaded
}

// Render will parse and compile
func (st *Stemplate) Render(w *http.ResponseWriter, templateName string) {

	if !st.LiveReload {
		st.templates[templateName].ExecuteTemplate(*w, "base", nil)
	} else {
		template.Must(template.ParseFiles(st.templatesDir+templateName, st.templatesDir+"base.tmpl")).ExecuteTemplate(*w, "base", nil)
	}

}
