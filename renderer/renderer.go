package renderer

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"vuk/config"
)

var (
	layoutMap map[string]([]byte)
	viewsMap  map[string]([]byte)
	Images    = map[string]string{}
	imgMux    = &sync.Mutex{}
)

func Init(layouts, views *[]string) (err error) {
	layoutMap = make(map[string][]byte, len(*layouts))
	viewsMap = make(map[string][]byte, len(*views))

	for _, layoutPath := range *layouts {
		rawLayout, err := os.ReadFile(layoutPath)
		if err != nil {
			return err
		}
		layoutMap[layoutPath] = rawLayout
	}

	for _, viewPath := range *views {
		rawView, err := os.ReadFile(viewPath)
		if err != nil {
			return err
		}
		viewsMap[viewPath] = rawView
	}

	return nil
}

func Render(page string) (handler func(w http.ResponseWriter, r *http.Request), err error) {
	rawPage, err := os.ReadFile(page)
	if err != nil {
		return nil, err
	}

	// TODO: Proper handling of layouts, as there won't be just one default layout.
	compiled := bytes.ReplaceAll(layoutMap[filepath.Join(config.DIR_LAYOUTS, "default.html")], []byte(MARK_PAGE), rawPage)

	viewsPaths := extractFilePaths(string(compiled))

	for len(viewsPaths) != 0 {
		for _, viewPath := range viewsPaths {
			normalizedViewPath := strings.TrimPrefix(viewPath, VIEW_START)
			normalizedViewPath = strings.TrimSuffix(normalizedViewPath, VIEW_END)
			normalizedViewPath = filepath.Join(config.DIR_VIEWS, normalizedViewPath)

			compiled = bytes.ReplaceAll(compiled, []byte(VIEW_START+viewPath+VIEW_END), viewsMap[normalizedViewPath])
		}

		extractedFilePaths := extractFilePaths(string(compiled))
		viewsPaths = []string{}
		for _, path := range extractedFilePaths {
			if IsImage(path) {
				imgMux.Lock()
				Images[path] = path
				imgMux.Unlock()
				continue
			}
			viewsPaths = append(viewsPaths, path)
		}
	}

	for realPath, renderedPath := range Images {
		compiled = bytes.ReplaceAll(compiled, []byte(VIEW_START+realPath+VIEW_END), []byte(renderedPath))
		continue
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(compiled)
	}, nil
}

func extractFilePaths(text string) (paths []string) {
	found := false

	token := ""
	for _, ch := range text {
		char := string(ch)
		token += char

		if strings.HasSuffix(token, VIEW_START) {
			found = true
			token = ""
			continue
		}
		if found && strings.HasSuffix(token, VIEW_END) {
			found = false
			paths = append(paths, strings.TrimSuffix(token, VIEW_END))
			token = ""
			continue
		}
	}

	return paths
}

var imageExtensions = []string{
	".png",
	".jpg",
	".jpeg",
}

func IsImage(path string) bool {
	for _, ext := range imageExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
