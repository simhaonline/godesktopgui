// A demo web server application, serving static content from a single endpoint.
// With the point of interest being that the html, css and javascript have
// been compiled-in to the executable.
package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/peterhoward42/godesktopgui/generate"
	"github.com/pkg/browser"
)

// htmlTemplate generates the HTML we serve to implement the GUI when we call
// its ExecuteTemplate method.
var htmlTemplate *template.Template

func main() {

	// Prepare the html template that will be combined with a data model to
	// serve html pages.

	htmlTemplate = parseTemplate()

	// The html we serve has href links to css and .js files - the URLs of which
	// start with /files, so we route all /files requests to the standard
	// library http.FileServer. The FileServer requires that we provide
	// an http.FileSystem. And that is how the compiled-in files present
	// themselves. See the generate package for how this gets created.

	http.Handle("/files/", http.FileServer(generate.CompiledFileSystem))

	// The GUI home page has its own dedicated handler.
	http.HandleFunc("/thegui", guiHandler)

	// Spin-up the standard library's http server in a separate goroutine.
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("http.ListenAndServe: %v", err)
		}
	}()

	// Give the server time to be ready.
	time.Sleep(3 * time.Second)

	// Then bring up a browser window or tab pointing to it.
	// Note this is asynchronous, and the call returns immediately.
	err := browser.OpenURL("http://127.0.0.1:8080/thegui")
	if err != nil {
		log.Fatalf("browser.Open: %v", err)
	}

	// Keep the main goroutine alive.
	wait := make(chan bool)
	<-wait

	log.Printf("Finished normally")

}

// parseTemplate retreives a template HTML file from the compiled-in
// file system, and parses it using the standard library Template.Parse
// to create a Template object.
func parseTemplate() *template.Template {
	fName := "files/templates/maingui.html"
	file, err := generate.CompiledFileSystem.Open(fName)
	if err != nil {
		log.Fatalf("Failed to open <%s>: %v", fName, err)
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read contents of file: %v", err)
	}
	t, err := template.New("gui").Parse(string(contents))
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
	return t
}

// GuiData holds the GUI state data that will be combined with the
// template to render the GUI HTML. Note it is a hierarchical structure -
// having a slice of TableRow - which the directives in the templating
// system is clever enough to iterate over.
type GuiData struct {
	Title       string
	Unwatch     int
	Star        int
	Fork        int
	Commits     int
	Branch      int
	Release     int
	Contributor int
	RowsInTable []TableRow
}

// TableRow is a slave model to describe a single row in an HTML table.
type TableRow struct {
	File    string
	Comment string
	Ago     string
	Icon    string
}

// guiHandler serves the GUI. Simple as that.
func guiHandler(w http.ResponseWriter, r *http.Request) {

	// Set the stateful parameters of the Gui data model - to create the
	// the dynamic user experience.
	dynamicData := populateGuiData()

	// This (standard library) call combines the template with the data model
	// to produce the required HTML. What is not obvious is that it does not
	// return the HTML here, but is capable, in of itself, of writing the HTML
	// it generates directly to the http.ResponseWriter provided.
	err := htmlTemplate.ExecuteTemplate(w, "gui", dynamicData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// populateGuiData populates a GuiData trivially with hard-coded data.
func populateGuiData() *GuiData {
	guiData := &GuiData{
		Title:       "Golang Standalone GUI Example",
		Unwatch:     3,
		Star:        0,
		Fork:        2,
		Commits:     31,
		Release:     1,
		Contributor: 1,
		RowsInTable: []TableRow{},
	}
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"do_this.go", "Initial commit", "1 month ago", "file"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"do_that.go", "Initial commit", "1 month ago", "file"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"index.go", "Initial commit", "1 month ago", "file"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"resources", "Initial commit", "2 months ago", "folder-open"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"docs", "Initial commit", "2 months ago", "folder-open"})
	return guiData
}
