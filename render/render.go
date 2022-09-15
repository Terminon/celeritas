package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Renderer struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{} // interface{} == empty interface, replaced by "any" keyword since Go 1.18
	// see also https://stackoverflow.com/questions/23148812/whats-the-meaning-of-interface
	// An interface value is constructed of two words of data:
	//one word is used to point to a method table for the value’s underlying type,
	//and the other word is used to point to the actual data being held by that value.
	CSRFToken  string
	Port       string
	ServerName string
	Secure     bool
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":

	}
	return nil
}

func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	// We need to take Parsefiles from html/template, NOT text/template
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)
	if err != nil {
		return err
	}

	return nil
}
