package train

import (
	"bytes"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func ReadAsset(assetUrl string) string {
	// TODO: buffer
	data := bytes.NewBuffer([]byte(""))

	fileExt := path.Ext(assetUrl)
	switch fileExt {
	case ".js":
		FindAssetsFunc(assetUrl, func(filePath string, content string) {
			data.Write([]byte(content + "\n"))
		})
	case ".css":
		FindAssetsFunc(assetUrl, func(filePath string, content string) {
			data.Write([]byte(content + "\n"))
		})
	case "":

	default:
		data.Write([]byte(ReadStaticAsset(assetUrl)))
	}

	// correct
	return string(data.Bytes())
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

func FindAssetsFunc(assetUrl string, found func(filePath string, content string)) (filePaths []string) {
	filePath := ResolvePath(assetUrl)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return filePaths
	}

	fileExt := path.Ext(assetUrl)
	head := string(patterns[fileExt]["head"].Find(content))

	if len(head) != 0 {
		content = patterns[fileExt]["head"].ReplaceAll(content, []byte(""))

		for _, line := range strings.Split(head, "\n") {
			// TODO: test match before reading files to avoid "test/*.css"
			assetUrl := patterns[fileExt]["require"].ReplaceAll([]byte(line), []byte(""))

			if len(assetUrl) == 0 {
				continue
			}

			paths := FindAssetsFunc(string(assetUrl)+fileExt, found)
			filePaths = append(filePaths, paths...)
		}
	}

	found(filePath, string(content))

	return append(filePaths, filePath)
}

func ResolvePath(assetUrl string) string {
	filePath := string(strings.Replace(assetUrl, Config.AssetsURL, "", 1))
	result := Config.AssetsPath + "/" + filePath
	result = string(regexp.MustCompile(`\/{2,}`).ReplaceAll([]byte(result), []byte("/")))

	return result
}

func ReadStaticAsset(assetUrl string) string {
	filePath := ResolvePath(assetUrl)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(content)
}
