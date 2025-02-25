package routing

import (
	"net/http"
	"strings"
	"sync"
	"vuk/config"
	"vuk/renderer"
)

func Init(r *http.ServeMux, pages *[]string) (commonErr error) {
	wg := &sync.WaitGroup{}
	mux := &sync.Mutex{}

	for _, page := range *pages {
		wg.Add(1)

		func() {
			defer wg.Done()

			handler, err := renderer.Render(page)
			if err != nil {
				commonErr = err
				return
			}

			normalizedPage := strings.TrimPrefix(page, config.DIR_PAGES)
			normalizedPage = strings.TrimPrefix(normalizedPage, "/")

			if normalizedPage == "index.html" {
				r.HandleFunc("GET /", handler)
				return
			}

			mux.Lock()
			defer mux.Unlock()
			r.HandleFunc("GET /"+strings.TrimSuffix(normalizedPage, ".html"), handler)
		}()
	}

	wg.Wait()

	return nil
}
