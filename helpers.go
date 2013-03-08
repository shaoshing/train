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

var HelperFuncs = template.FuncMap{
	"javascript_tag":            JavascriptTag,
	"stylesheet_tag":            StylesheetTag,
	"stylesheet_tag_with_param": StylesheetTagWithParam,
}

const (
	javascriptTag = `<script src="%s"%s></script>`
	stylesheetTag = `<link type="text/css" rel="stylesheet" href="%s"%s>`
)

func JavascriptTag(name string) template.HTML {
	assetUrl := "javascripts/" + name + ".js"
	paths, mtimes := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, "", mtimes, javascriptTag)
}

func StylesheetTagWithParam(name string, param string) template.HTML {
	assetUrl := "stylesheets/" + name + ".css"
	paths, mtimes := resolveAssetUrls(assetUrl)
	return generateRawHtml(paths, param, mtimes, stylesheetTag)
}

func StylesheetTag(name string) template.HTML {
	return StylesheetTagWithParam(name, "")
}

func resolveAssetUrls(assetUrl string) (urls []string, mtimes []time.Time) {
	if IsInProduction() {
		return getBundledAssets(assetUrl)
	}
	return getUnbundledAssets(assetUrl)
}

func getUnbundledAssets(assetUrl string) (urls []string, mtimes []time.Time) {
	filePath := ResolvePath(assetUrl)
	fileExt := path.Ext(filePath)
	var paths []string

	if fileExt == ".js" || fileExt == ".css" {
		var err error
		paths, err = ReadAssetsFunc(filePath, assetUrl, func(filePath string, content string) {})
		if err != nil {
			panic(err)
		}
	} else {
		paths = append(paths, filePath)
	}

	for _, assetPath := range paths {
		info, _ := os.Stat(assetPath)
		mtimes = append(mtimes, info.ModTime())
		urls = append(urls, asserUrlFromPath(assetPath))
	}
	return
}

func getBundledAssets(assetUrl string) (urls []string, mtimes []time.Time) {
	urls = []string{ManifestInfo[Config.AssetsUrl+assetUrl]}
	mtimes = nil

	return
}

func asserUrlFromPath(assetPath string) (url string) {
	url = path.Clean(strings.Replace(assetPath, Config.AssetsPath, Config.AssetsUrl, 1))
	url = strings.Replace(url, ".sass", ".css", 1)
	url = strings.Replace(url, ".scss", ".css", 1)
	url = strings.Replace(url, ".coffee", ".js", 1)

	return
}

func generateRawHtml(urls []string, param string, mtimes []time.Time, tag string) template.HTML {
	htmls := []string{}
	if len(param) != 0 {
		param = " " + param
	}

	for i, url := range urls {
		murl := url
		if mtimes != nil {
			murl += "?" + strconv.FormatInt(mtimes[i].Unix(), 10)
		}
		htmls = append(htmls, fmt.Sprintf(tag, murl, param))
	}
	return template.HTML(strings.Join(htmls, "\n"))
}
