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
	data         map[string]interface{}
}

//NewStemplate create a new Istance of Stemplate obj
func NewStemplate(templatesDirectory string) (*Stemplate, error) {
	var st Stemplate

	st.templatesDir = templatesDirectory
	st.LiveReload = false
	st.templates = make(map[string]*template.Template)
	st.data = make(map[string]interface{})

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

// loadTemplate will load at run time the desidered temp
func (st *Stemplate) loadTemplate(tname string) *template.Template {

	templates, terr := filepath.Glob(st.templatesDir + "*.tmpl")
	if terr != nil {
		fmt.Println("[JIT template]: ERROR ~ " + terr.Error())
		return nil
	}

	templates = append(templates, st.templatesDir+tname)

	return template.Must(template.ParseFiles(templates...))
}

// Reload is an utily function that allow to recompile templates at run time
func (st *Stemplate) Reload() {

	// useless if we are using live reload
	if st.LiveReload {
		return
	}

	stReloaded, err := NewStemplate(st.templatesDir)
	if err != nil {
		fmt.Println("[Reloading templates]: ERROR")
		return
	}

	stReloaded.data = st.data

	st = stReloaded
}

// Render will parse and compile
func (st *Stemplate) Render(w *http.ResponseWriter, templateName string) {

	if !st.LiveReload {
		st.templates[templateName].ExecuteTemplate(*w, "base", st.data)
	} else {
		st.loadTemplate(templateName).ExecuteTemplate(*w, "base", st.data)
	}

}

// AddData will add the desidered obj to a template and will add it every time
// that template is executed
func (st *Stemplate) AddData(name string, data interface{}) {
	st.data[name] = data
}

// RemoveData is used when you dont need specific obj to be passed anymore to your template
func (st *Stemplate) RemoveData(name string) {
	delete(st.data, name)
}
