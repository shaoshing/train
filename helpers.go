package train

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Helpers struct{}

const (
	JavascriptTag = `<script src="%s"></script>`
	StylesheetTag = `<link type="text/css" rel="stylesheet" href="%s">`
)

func (this Helpers) JavascriptTag(name string) template.HTML {
	assetUrl := "javascripts/" + name + ".js"
	paths, mtimes := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, mtimes, JavascriptTag)
}

func (this Helpers) StylesheetTag(name string) template.HTML {
	assetUrl := "stylesheets/" + name + ".css"
	paths, mtimes := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, mtimes, StylesheetTag)
}

func resolveAssetUrls(assetUrl string) (urls []string, mtimes []time.Time) {
	paths := FindAssetsFunc(assetUrl, func(filePath string, content string) {})
	if Config.BundleAssets {
		paths = paths[len(paths)-1:]
	}

	for _, assetPath := range paths {
		info, _ := os.Stat(assetPath)
		mtimes = append(mtimes, info.ModTime())
		assetUrl := path.Clean(strings.Replace(assetPath, Config.AssetsPath, Config.AssetsUrl, 1))
		urls = append(urls, assetUrl)
	}
	return
}

func generateRawHtml(urls []string, mtimes []time.Time, tag string) template.HTML {
	htmls := []string{}
	for i, url := range urls {
		murl := url + "?" + strconv.FormatInt(mtimes[i].Unix(), 10)
		htmls = append(htmls, fmt.Sprintf(tag, murl))
	}
	return template.HTML(strings.Join(htmls, "\n"))
}
