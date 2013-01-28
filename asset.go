package train

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

func ReadAsset(assetUrl string) string {
	// TODO: buffer
	data := ""

	fileExt := path.Ext(assetUrl)
	switch fileExt {
	case ".js":
		FindAssetsFunc(assetUrl, func(filePath string, content string) {
			data += content
		})
	case ".css":
		FindAssetsFunc(assetUrl, func(filePath string, content string) {
			data += content
		})
	case "":

	default:
		data = ReadStaticAsset(assetUrl)
	}

	return string(data)
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
	fileExt := path.Ext(assetUrl)
	var assetFolder string
	switch fileExt {
	case ".js":
		assetFolder = "/javascripts/"
	case ".css":
		assetFolder = "/stylesheets/"
	}

	filePath := string(regexp.MustCompile(`\/{2,}`).ReplaceAll([]byte(strings.Replace(assetUrl, Config.AssetsURL, "", 1)), []byte("/")))
	return Config.AssetsPath + assetFolder + filePath
}

func ReadStaticAsset(assetUrl string) string {
	assetUrl = strings.Replace(assetUrl, Config.AssetsURL, "", 1)
	for _, assetPath := range []string{"/", "/javascripts/", "/stylesheets/", "/images/"} {
		filePath := Config.AssetsPath + assetPath + assetUrl
		filePath = string(regexp.MustCompile(`\/{2,}`).ReplaceAll([]byte(filePath), []byte("/")))
		content, err := ioutil.ReadFile(filePath)
		if err == nil {
			return string(content)
		}
	}
	return ""
}
