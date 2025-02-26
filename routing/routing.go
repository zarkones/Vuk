package routing

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

	for realPath, renderredPath := range renderer.Images {
		r.HandleFunc("GET /"+renderredPath, func() func(w http.ResponseWriter, r *http.Request) {
			imgPath := filepath.Join(config.DIR_VIEWS, realPath)
			image, err := os.ReadFile(imgPath)
			if err != nil {
				fmt.Println("failed to read image at:", imgPath)
			}
			return func(w http.ResponseWriter, r *http.Request) {
				w.Write(image)
			}
		}())
	}

	return nil
}
