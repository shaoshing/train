package train

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shaoshing/train/interpreter"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func ReadAsset(assetUrl string) (result string, err error) {
	filePath := ResolvePath(assetUrl)
	fileExt := path.Ext(filePath)

	switch fileExt {
	case ".js", ".css", ".scss", ".sass", ".coffee":
		if Config.BundleAssets {
			data := bytes.NewBuffer([]byte(""))
			contents := []string{}
			_, err = ReadAssetsFunc(filePath, assetUrl, func(filePath string, content string) {
				contents = append(contents, content)
			})
			if err != nil {
				if !IsInProduction() {
					fmt.Println(err)
				}
				return
			}
			data.Write([]byte(strings.Join(contents, "\n")))
			result = string(data.Bytes())
		} else {
			result, err = ReadRawAndCompileAsset(filePath, assetUrl)
		}
	default:
		err = errors.New("Unsupported Asset: " + assetUrl)
	}

	return
}

func compileSassAndCoffee(filePath string) (string, error) {
	interpreter.Config.AssetsPath = Config.AssetsPath
	interpreter.Config.SASS.LineNumbers = Config.SASS.LineNumbers
	interpreter.Config.SASS.DebugInfo = Config.SASS.DebugInfo
	interpreter.Config.Verbose = Config.Verbose
	result, err := interpreter.Compile(filePath)
	return result, err
}

var patterns = map[string](map[string]*regexp.Regexp){
	".js": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(.*\n)*(\ *\/\/\=\ *require\ +.*\n)+`),
		"require": regexp.MustCompile(`^\ *\/\/\=\ *require\ +`),
	},
	".coffee": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(.*\n)*(#=\ *require\ +.*\n)+`),
		"require": regexp.MustCompile(`^#=\ *require\ +`),
	},
	".css": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(.*\n)*(\ *\*\=\ *require\ +.*\n)+(\ *\*\/\ *\n)`),
		"require": regexp.MustCompile(`^\ *\*\=\ *require\ +`),
	},
}

func patternExt(fileExt string) string {
	switch fileExt {
	case ".scss", ".sass":
		return ".css"
	}

	return fileExt
}

func ReadAssetsFunc(filePath, assetUrl string, found func(filePath string, content string)) (filePaths []string, err error) {
	fileExt := path.Ext(filePath)
	fileExtPattern := patternExt(fileExt)
	var cacheKey string
	if cacheKey, err = generateCacheKey(filePath); err != nil {
		err = errors.New("Asset Not Found: " + assetUrl)
		return
	}

	requiredAssetUrls, content, hit := readFromCache(cacheKey)
	if !hit {
		content, err = ReadRawAsset(filePath, assetUrl)
		if err != nil {
			return
		}
		header := FindDirectivesHeader(&content, fileExtPattern)
		content, err = ReadRawAndCompileAsset(filePath, assetUrl)
		content = hotFixSASSCommentLines(content, fileExt)
		if len(header) != 0 {
			content = strings.Replace(content, header, "", 1)

			for _, line := range strings.Split(header, "\n") {
				if !patterns[fileExtPattern]["require"].Match([]byte(line)) {
					continue
				}

				requiredAssetUrl := string(patterns[fileExtPattern]["require"].ReplaceAll([]byte(line), []byte("")))
				if len(requiredAssetUrl) == 0 {
					continue
				}

				requiredAssetUrls = append(requiredAssetUrls, requiredAssetUrl)
			}
		}
		writeToCache(cacheKey, requiredAssetUrls, content)
	}

	for _, requiredAssetUrl := range requiredAssetUrls {
		var paths []string
		requiredFilePath := ResolvePath(requiredAssetUrl + fileExt)
		requiredFileExt := path.Ext(requiredFilePath)
		paths, err = ReadAssetsFunc(requiredFilePath, requiredAssetUrl+requiredFileExt, found)
		if err != nil {
			err = errors.New(fmt.Sprintf("%s\n--- required by %s", err.Error(), assetUrl))
			return
		}

		filePaths = append(filePaths, paths...)
	}

	found(filePath, content)
	filePaths = append(filePaths, filePath)
	return
}

func FindDirectivesHeader(content *string, fileExt string) string {
	return string(patterns[fileExt]["head"].Find([]byte(*content)))
}

var mapAlterExtensions = map[string]string{
	".css":    ".sass|.scss",
	".sass":   ".scss|.css",
	".scss":   ".sass|.css",
	".js":     ".coffee",
	".coffee": ".js",
}

// Find possible asset files.
// url = javascripts/asset.js
// if url exist
// => javascript/asset.js
// or url alternation exist
// => javascript/asset.coffee
func ResolvePath(assetUrl string) (assetPath string) {
	assetPath = string(strings.Replace(assetUrl, Config.AssetsUrl, "", 1))
	assetPath = path.Clean(Config.AssetsPath + "/" + assetPath)
	fileExt := path.Ext(assetPath)
	alterExts, hasAlterExt := mapAlterExtensions[fileExt]
	if !isFileExist(assetPath) && hasAlterExt {
		for _, alterExt := range strings.Split(alterExts, "|") {
			alterPath := strings.Replace(assetPath, fileExt, alterExt, 1)
			if isFileExist(alterPath) {
				assetPath = alterPath
				break
			}
		}
	}

	return
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Hotfix SASS commet lines
// from:
//    /*= require foo
// to:
//    /*
//     *= require foo
func hotFixSASSCommentLines(content string, fileExt string) string {
	if fileExt == ".sass" {
		return strings.Replace(content, "/*= require", "/*\n *= require", 1)
	}

	return content
}

func ReadRawAsset(filePath, assetUrl string) (result string, err error) {
	content, _err := ioutil.ReadFile(filePath)
	if _err != nil {
		err = errors.New("Asset Not Found: " + assetUrl)
		return
	}
	result = string(content)

	return
}

// .js, .css read raw
// .coffee, .scss, .sass will compile
func ReadRawAndCompileAsset(filePath, assetUrl string) (result string, err error) {
	fileExt := path.Ext(filePath)

	if fileExt == ".scss" || fileExt == ".sass" || fileExt == ".coffee" {
		result, err = compileSassAndCoffee(filePath)
	} else {
		result, err = ReadRawAsset(filePath, assetUrl)
	}

	return
}

type AssetCache struct {
	RequiredUrls []string
	Content      string
}

var assetsCache = make(map[string]AssetCache)

func generateCacheKey(filePath string) (key string, err error) {
	var info os.FileInfo
	info, err = os.Stat(filePath)
	if err != nil {
		return
	}
	key = filePath + strconv.FormatInt(info.ModTime().Unix(), 10)
	return
}

func readFromCache(key string) (requiredUrls []string, content string, hit bool) {
	var cache AssetCache
	if cache, hit = assetsCache[key]; hit {
		requiredUrls = cache.RequiredUrls
		content = cache.Content
	}
	return
}

func writeToCache(key string, requiredUrls []string, content string) {
	assetsCache[key] = AssetCache{requiredUrls, content}
}
