package vuk

import (
	"net/http"
	"vuk/config"
	"vuk/helpers"
	"vuk/renderer"
	"vuk/routing"
)

var (
	layouts []string
	pages   []string
	views   []string
)

func Init(r *http.ServeMux) (err error) {
	layouts, err = helpers.WalkDir(config.DIR_LAYOUTS)
	if err != nil {
		return err
	}
	pages, err = helpers.WalkDir(config.DIR_PAGES)
	if err != nil {
		return err
	}
	views, err = helpers.WalkDir(config.DIR_VIEWS)
	if err != nil {
		return err
	}

	if err := renderer.Init(&layouts, &views); err != nil {
		return err
	}

	if err := routing.Init(r, &pages); err != nil {
		return err
	}

	renderer.ExportRenderedItems()

	return nil
}
