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
		interpreter.Config.SASS.LineNumbers = Config.SASS.LineNumbers
		interpreter.Config.SASS.DebugInfo = Config.SASS.DebugInfo
		result, err = interpreter.Compile(filePath)
	default:
		err = errors.New("Unsupported Asset: " + assetUrl)
	}
	return

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
	var content string
	content, err = ReadRawAsset(filePath, assetUrl)
	if err != nil {
		return
	}

	fileExt := path.Ext(filePath)
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

			var paths []string
			requiredFilePath := ResolvePath(requiredAssetUrl + fileExt)
			paths, err = ReadAssetsFunc(requiredFilePath, requiredAssetUrl+fileExt, found)
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

var mapAlterExtensions = map[string]string{
	".css": ".sass|.scss",
	".js":  ".coffee",
}

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
