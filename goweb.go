package web

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/deryl-sagala/logger"
)

func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var formattedContent string
	for scanner.Scan() {
		formattedContent += fmt.Sprintf("\t\t%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return formattedContent, nil
}

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
func RenderHTML(tmpl string, vars ...map[string]interface{}) http.HandlerFunc {
	log := logger.NewLogger()
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := readFile(fmt.Sprintf("templates/%s", tmpl))
		if err != nil {
			log.Error(fmt.Sprintf("", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(vars) > 0 && vars[0] != nil {
			for key, value := range vars[0] {
				placeholder := "{{" + key + "}}"
				result = strings.Replace(result, placeholder, fmt.Sprintf("%v", value), -1)
			}
		}

		// Render the modified template to the response
		_, err = w.Write([]byte(result))
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func Return(tmpl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger()
		// Set the Content-Type header to specify that the response contains plain text.
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(tmpl))
		if err != nil {
			log.Error(err.Error())
			// If an error occurs while writing, return an internal server error.
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
// Deprecated: Use renderer directly
func Wrap(h func()) http.HandlerFunc {
	log := logger.NewLogger()
	log.Warn("Wrap function is deprecated and does nothing and is kept for backward compatibility.")
	// Return an empty http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		// This function does nothing and is kept for backward compatibility.
		// You can choose to log or notify about the deprecation here if you want.
	}
}
