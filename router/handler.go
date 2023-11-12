package router

import (
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/abspayd/music-companion/music"
)

var (
	validPath = regexp.MustCompile("^/(home|intervals)$")
	templates = loadTemplates()
)

func loadTemplates() *template.Template {
	tmplFS := os.DirFS("./tmpl")
	return template.Must(template.ParseFS(tmplFS, "*.html"))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[1])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		// Redirect "/" to "/home"
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		// This url was not matched by any handlers
		http.NotFound(w, r)
		return
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request) {

}

func handleIndex(w http.ResponseWriter, r *http.Request, tmpl string) {
	renderTemplate(w, tmpl+".html", nil)
}

func handleIntervals(w http.ResponseWriter, r *http.Request, tmpl string) {
	if r.Method == http.MethodGet {
		renderTemplate(w, tmpl+".html", nil)
	} else if r.Method == http.MethodPost {
		p1 := r.FormValue("pitch1")
		p2 := r.FormValue("pitch2")

		idx1, _ := music.Search(p1)
		idx2, _ := music.Search(p2)

		note1 := music.Note{Pitch: idx1, Octave: 0}
		note2 := music.Note{Pitch: idx2, Octave: 0}

		interval := note1.GetInterval(note2)
		intervalString := music.IntervalToString(interval)

		templates.ExecuteTemplate(w, "intervalResult", intervalString)
	}
}

func handleValidateNote(w http.ResponseWriter, r *http.Request) {
	pitch := "c" // r.PostForm
	_, err := music.Search(pitch)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	// w.Write([]byte("Validating input..."))
}