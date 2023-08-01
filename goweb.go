package web

import (
	"html/template"
	"net/http"
)

var templates *template.Template

type route struct {
	pattern string
	handler http.HandlerFunc
}

var routes []route

/*
	Function to register a new route

example: web.AddRoute("/", web.wrap(index))
example wrapper function:

	func index() {
		web.RenderHTML(w, "index.html")
	}
*/
func AddRoute(pattern string, handler http.HandlerFunc) {
	routes = append(routes, route{pattern: pattern, handler: handler})
}
func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

/*
web.renderHtml(file.html) to render html template
make sure you have file.html in /templates/
*/
func renderHTML(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func returnText(w http.ResponseWriter, text string) {
	// Set the Content-Type header to specify that the response contains plain text.
	w.Header().Set("Content-Type", "text/plain")

	// Write the text to the ResponseWriter.
	_, err := w.Write([]byte(text))
	if err != nil {
		// If an error occurs while writing, return an internal server error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// web.serve(port) to serve your website on the specified port
func Serve(port string) {
	for _, route := range routes {
		http.HandleFunc(route.pattern, route.handler)
	}
	http.ListenAndServe(":"+port, nil)
}

// A helper function to wrap the handler functions without explicit w and r parameters
func wrap(h func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h()
	}
}
