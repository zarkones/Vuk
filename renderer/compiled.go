package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	EXPORT_DIR = "export"
)

var compiledItem = map[string][]byte{}

var compiledItemMux = &sync.Mutex{}

func AddRenderedItem(path string, content []byte) {
	defer compiledItemMux.Unlock()
	compiledItemMux.Lock()
	compiledItem[path] = content
}

func ExportRenderedItems() (err error) {
	os.MkdirAll(EXPORT_DIR, 0777)

	for endpointPath, content := range compiledItem {
		_ = content
		endpointPath = strings.TrimPrefix(endpointPath, "/")
		if len(endpointPath) == 0 {
			endpointPath = "index.html"
		}
		if !strings.Contains(endpointPath, ".") {
			endpointPath += ".html"
		}

		dir := filepath.Dir(endpointPath)

		if dir != "." {
			os.MkdirAll(filepath.Join(EXPORT_DIR, dir), 0777)
		}

		os.WriteFile(filepath.Join(EXPORT_DIR, endpointPath), content, 0777)
	}
	return nil
}
