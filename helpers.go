package train

import (
	"html/template"
	"regexp"
	"strings"
)

type Helpers struct{}

const (
	JavascriptTemplate = `<script src="{path}"></script>`
	StylesheetTemplate = `<link type="text/css" rel="stylesheet" href="{path}">`
)

func (this Helpers) JavascriptIncludeTag(name string) template.HTML {
	assetUrl := "javascripts/" + name + ".js"
	paths := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, JavascriptTemplate)
}

func (this Helpers) StylesheetIncludeTag(name string) template.HTML {
	assetUrl := "stylesheets/" + name + ".css"
	paths := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, StylesheetTemplate)
}

func resolveAssetUrls(assetUrl string) []string {
	paths := FindAssetsFunc(assetUrl, func(filePath string, content string) {})
	if Config.BundleAssets {
		paths = paths[len(paths)-1:]
	}

	urls := make([]string, len(paths))
	for i, path := range paths {
		assetUrl := strings.Replace(path, Config.AssetsPath, Config.AssetsUrl, 1)
		assetUrl = string(regexp.MustCompile(`\/{2,}`).ReplaceAll([]byte(assetUrl), []byte("/")))
		urls[i] = assetUrl
	}
	return urls
}

var pathReg = regexp.MustCompile(`\{path\}`)

func generateRawHtml(paths []string, html string) template.HTML {
	result := ""
	for i, path := range paths {
		result += string(pathReg.ReplaceAll([]byte(html), []byte(path)))
		if i != len(paths)-1 {
			result += "\n"
		}
	}
	return template.HTML(result)
}
