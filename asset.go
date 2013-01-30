package train

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func ReadAsset(assetUrl string) (result string, err error) {
	fileExt := path.Ext(assetUrl)
	if fileExt != ".js" && fileExt != ".css" {
		err = errors.New("Unsupported Asset: " + assetUrl)
		return
	}

	if Config.BundleAssets {
		data := bytes.NewBuffer([]byte(""))
		contents := []string{}
		_, err = ReadAssetsFunc(assetUrl, func(filePath string, content string) {
			contents = append(contents, content)
		})
		if err != nil {
			return
		}
		data.Write([]byte(strings.Join(contents, "\n")))
		result = string(data.Bytes())
	} else {
		result, err = ReadRawAsset(assetUrl)
	}
	return
}

var patterns = map[string](map[string]*regexp.Regexp){
	".js": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(\/\/\=\ require\ +.*\n)+`),
		"require": regexp.MustCompile(`^\/\/\=\ require\ +`),
	},
	".css": map[string]*regexp.Regexp{
		"head":    regexp.MustCompile(`(\ *\/\*\ *\n)(\ *\*\=\ +require\ +.*\n)+(\ *\*\/\ *\n)`),
		"require": regexp.MustCompile(`^\ *\*\=\ +require\ +`),
	},
}

func ReadAssetsFunc(assetUrl string, found func(filePath string, content string)) (filePaths []string, err error) {
	filePath := ResolvePath(assetUrl)

	var content string
	content, err = ReadRawAsset(assetUrl)
	if err != nil {
		return
	}

	fileExt := path.Ext(assetUrl)
	header := FindDirectivesHeader(&content, fileExt)

	if len(header) != 0 {
		content = strings.Replace(content, header, "", 1)

		for _, line := range strings.Split(header, "\n") {
			if !patterns[fileExt]["require"].Match([]byte(line)) {
				continue
			}

			requiredAssetUrl := patterns[fileExt]["require"].ReplaceAll([]byte(line), []byte(""))
			if len(requiredAssetUrl) == 0 {
				continue
			}

			var paths []string
			paths, err = ReadAssetsFunc(string(requiredAssetUrl)+fileExt, found)
			if err != nil {
				err = errors.New(fmt.Sprintf("%s\n--- required by %s", err.Error(), assetUrl))
				return
			}

			filePaths = append(filePaths, paths...)
		}
	}

	found(filePath, content)
	filePaths = append(filePaths, filePath)
	return
}

func FindDirectivesHeader(content *string, fileExt string) string {
	return string(patterns[fileExt]["head"].Find([]byte(*content)))
}

func ResolvePath(assetUrl string) string {
	filePath := string(strings.Replace(assetUrl, Config.AssetsUrl, "", 1))
	result := path.Clean(Config.AssetsPath + "/" + filePath)

	return result
}

func ReadRawAsset(assetUrl string) (result string, err error) {
	filePath := ResolvePath(assetUrl)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errors.New("Asset Not Found: " + assetUrl)
		return
	}
	result = string(content)

	return
}
