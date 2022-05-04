package server

import (
	"embed"
)

var (
	//go:embed template/layouts/* template/blocks/*
	templates embed.FS
	//go:embed assets/dist/*
	assets embed.FS
)

func getEmbeddedTemplates() embed.FS {
	return templates
}

func getEmbeddedAssets() embed.FS {
	return assets
}
