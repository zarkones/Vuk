package routing

import (
	"net/http"
	"strings"
	"vuk/config"
	"vuk/renderer"
)

func Init(r *http.ServeMux, pages *[]string) (err error) {
	for _, page := range *pages {
		handler, err := renderer.Render(page)
		if err != nil {
			return err
		}

		normalizedPage := strings.TrimPrefix(page, config.DIR_PAGES)
		normalizedPage = strings.TrimPrefix(normalizedPage, "/")

		if normalizedPage == "index.html" {
			r.HandleFunc("GET /", handler)
			continue
		}
		r.HandleFunc("GET /"+strings.TrimSuffix(normalizedPage, ".html"), handler)
	}
	return nil
}
