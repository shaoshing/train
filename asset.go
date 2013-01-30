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

	if Config.BundleAssets {
		fileExt := path.Ext(assetUrl)
		switch fileExt {
		case ".js", ".css":
			contents := []string{}
			FindAssetsFunc(assetUrl, func(filePath string, content string) {
				contents = append(contents, content)
			})
			data.Write([]byte(strings.Join(contents, "\n")))
		case "":

		default:
			data.Write([]byte(ReadStaticAsset(assetUrl)))
		}
	} else {
		data.Write([]byte(ReadStaticAsset(assetUrl)))
	}
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

	b_content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return filePaths
	}
	content := string(b_content)

	fileExt := path.Ext(assetUrl)
	header := FindDirectivesHeader(&content, fileExt)

	if len(header) != 0 {
		content = strings.Replace(content, header, "", 1)

		for _, line := range strings.Split(header, "\n") {
			// TODO: test match before reading files to avoid "test/*.css"
			assetUrl := patterns[fileExt]["require"].ReplaceAll([]byte(line), []byte(""))

			if len(assetUrl) == 0 {
				continue
			}

			paths := FindAssetsFunc(string(assetUrl)+fileExt, found)
			filePaths = append(filePaths, paths...)
		}
	}

	found(filePath, content)

	return append(filePaths, filePath)
}

func FindDirectivesHeader(content *string, fileExt string) string {
	return string(patterns[fileExt]["head"].Find([]byte(*content)))
}

func ResolvePath(assetUrl string) string {
	filePath := string(strings.Replace(assetUrl, Config.AssetsUrl, "", 1))
	result := path.Clean(Config.AssetsPath + "/" + filePath)

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
