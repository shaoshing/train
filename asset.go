package train

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func ReadAsset(assetUrl string) (result string, err error) {
	fileExt := path.Ext(assetUrl)
	if fileExt != ".js" && fileExt != ".css" {
		err = errors.New("Can only read from js and css assets.")
		return
	}

	if Config.BundleAssets {
		data := bytes.NewBuffer([]byte(""))
		contents := []string{}
		_, err = FindAssetsFunc(assetUrl, func(filePath string, content string) {
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

func FindAssetsFunc(assetUrl string, found func(filePath string, content string)) (filePaths []string, err error) {
	filePath := ResolvePath(assetUrl)

	var b_content []byte
	b_content, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	content := string(b_content)

	fileExt := path.Ext(assetUrl)
	header := FindDirectivesHeader(&content, fileExt)

	if len(header) != 0 {
		content = strings.Replace(content, header, "", 1)

		for _, line := range strings.Split(header, "\n") {
			if !patterns[fileExt]["require"].Match([]byte(line)) {
				continue
			}

			assetUrl := patterns[fileExt]["require"].ReplaceAll([]byte(line), []byte(""))

			if len(assetUrl) == 0 {
				continue
			}

			var paths []string
			paths, err = FindAssetsFunc(string(assetUrl)+fileExt, found)
			if err != nil {
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
	if err == nil {
		result = string(content)
	}

	return
}
