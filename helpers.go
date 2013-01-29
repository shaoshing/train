package train

import (
	"html/template"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Helpers struct{}

const (
	JavascriptTemplate = `<script src="{path}"></script>`
	StylesheetTemplate = `<link type="text/css" rel="stylesheet" href="{path}">`
)

func (this Helpers) JavascriptIncludeTag(name string) template.HTML {
	assetUrl := "javascripts/" + name + ".js"
	paths, mtimes := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, mtimes, JavascriptTemplate)
}

func (this Helpers) StylesheetIncludeTag(name string) template.HTML {
	assetUrl := "stylesheets/" + name + ".css"
	paths, mtimes := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, mtimes, StylesheetTemplate)
}

func resolveAssetUrls(assetUrl string) (urls []string, mtimes []time.Time) {
	paths := FindAssetsFunc(assetUrl, func(filePath string, content string) {})
	if Config.BundleAssets {
		paths = paths[len(paths)-1:]
	}

	for _, path := range paths {
		info, _ := os.Stat(path)
		mtimes = append(mtimes, info.ModTime())
		assetUrl := strings.Replace(path, Config.AssetsPath, Config.AssetsUrl, 1)
		assetUrl = string(regexp.MustCompile(`\/{2,}`).ReplaceAll([]byte(assetUrl), []byte("/")))
		urls = append(urls, assetUrl)
	}
	return
}

var pathReg = regexp.MustCompile(`\{path\}`)

func generateRawHtml(urls []string, mtimes []time.Time, html string) template.HTML {
	htmls := []string{}
	for i, url := range urls {
		murl := url + "?" + strconv.FormatInt(mtimes[i].Unix(), 10)
		htmls = append(htmls, string(pathReg.ReplaceAll([]byte(html), []byte(murl))))
	}
	return template.HTML(strings.Join(htmls, "\n"))
}
