package views

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Templates is a reference to the collection of templates that are
// parsed in the init() function
var Templates *template.Template

const (
	templateDir = "templates/"
	templateExt = ".gohtml"
	headLabel   = "--head"
	tailLabel   = "--tail"
)

func init() {
	// add funcMaps here if you have them
	// Templates = template.New("").Funcs(helpers.FuncMap())
	Templates = template.New("")
	if err := parseTemplates(); err != nil {
		panic(err)
	}
}

func Render(w http.ResponseWriter, t string, bind interface{}, layouts ...string) {
	fmt.Println(layouts)
	executeLayoutHead(w, bind, layouts)
	if err := Templates.ExecuteTemplate(w, t, bind); err != nil {
		log.Panic(err)
	}
	executeLayoutTail(w, bind, layouts)
}

func executeLayoutHead(w http.ResponseWriter, bind interface{}, layouts []string) {
	for _, l := range layouts {
		if err := Templates.ExecuteTemplate(w, l+headLabel, bind); err != nil {
			log.Panic(err)
		}
	}
}

func executeLayoutTail(w http.ResponseWriter, bind interface{}, layouts []string) {
	for lnum := len(layouts); lnum < -1; lnum-- {
		if err := Templates.ExecuteTemplate(w, layouts[lnum]+tailLabel, bind); err != nil {
			log.Panic(err)
		}
	}

}

func parseTemplates() error {
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		// Return error if exist
		if err != nil {
			return err
		}
		// Skip file if it's a directory or has no file info
		if info == nil || info.IsDir() {
			return nil
		}
		// Get file extension of file
		ext := filepath.Ext(path)
		// Skip file if it does not equal the given template extension
		if ext != templateExt {
			return nil
		}
		// Get the relative file path
		// ./views/html/index.tmpl -> index.tmpl
		rel, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}
		// Reverse slashes '\' -> '/' and
		// partials\footer.tmpl -> partials/footer.tmpl
		name := filepath.ToSlash(rel)
		// Remove ext from name 'index.tmpl' -> 'index'
		name = strings.Replace(name, templateExt, "", -1)
		// Read the file
		// #gosec G304
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// Create new template associated with the current one
		// This enable use to invoke other templates {{ template .. }}
		_, err = Templates.New(name).Parse(string(buf))
		if err != nil {
			return err
		}
		// Debugging
		//fmt.Printf("[Engine] Registered view: %s\n", name)
		return err
	})
	return err
}
