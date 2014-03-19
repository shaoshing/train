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
	case ".js", ".css":
		if Config.BundleAssets {
			data := bytes.NewBuffer([]byte(""))
			contents := []string{}
			_, err = ReadAssetsFunc(filePath, assetUrl, func(filePath string, content string) {
				contents = append(contents, content)
			})
			if err != nil {
				return
			}
			data.Write([]byte(strings.Join(contents, "\n")))
			result = string(data.Bytes())
		} else {
			result, err = ReadRawAsset(filePath, assetUrl)
		}
	case ".sass", ".scss", ".coffee":
		interpreter.Config.AssetsPath = Config.AssetsPath
		interpreter.Config.SASS.LineNumbers = Config.SASS.LineNumbers
		interpreter.Config.SASS.DebugInfo = Config.SASS.DebugInfo
		interpreter.Config.Verbose = Config.Verbose
		result, err = interpreter.Compile(filePath)
	default:
		err = errors.New("Unsupported Asset: " + assetUrl)
	}

	return
}

var patterns = map[string](map[string]*regexp.Regexp){
	".js": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(.*\n)*(\ *\/\/\=\ *require\ +.*\n)+`),
		"require": regexp.MustCompile(`^\ *\/\/\=\ *require\ +`),
	},
	".css": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(.*\n)*(\ *\/\*\ *\n)(\ *\*\=\ *require\ +.*\n)+(\ *\*\/\ *\n)`),
		"require": regexp.MustCompile(`^\ *\*\=\ *require\ +`),
	},
}

func ReadAssetsFunc(filePath, assetUrl string, found func(filePath string, content string)) (filePaths []string, err error) {
	fileExt := path.Ext(filePath)
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
		header := FindDirectivesHeader(&content, fileExt)

		if len(header) != 0 {
			content = strings.Replace(content, header, "", 1)

			for _, line := range strings.Split(header, "\n") {
				if !patterns[fileExt]["require"].Match([]byte(line)) {
					continue
				}

				requiredAssetUrl := string(patterns[fileExt]["require"].ReplaceAll([]byte(line), []byte("")))
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
		paths, err = ReadAssetsFunc(requiredFilePath, requiredAssetUrl+fileExt, found)
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
	".css": ".sass|.scss",
	".js":  ".coffee",
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

func ReadRawAsset(filePath, assetUrl string) (result string, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("Asset Not Found: " + assetUrl)
		return
	}
	result = string(content)

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
