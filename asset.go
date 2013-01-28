package train

import (
	"io/ioutil"
	"path"
	"strings"
)

func ReadAsset(url string) string {
	filePath := strings.Replace(url, Config.AssetsURL, "/", 1)
	fileExt := path.Ext(url)

	var data []byte
	var err error
	switch fileExt {
	case ".js":
		data, err = ioutil.ReadFile(Config.AssetsPath + "/javascripts" + filePath)
		if err != nil {
			return ""
		}
	case ".css":
		data, err = ioutil.ReadFile(Config.AssetsPath + "/stylesheets" + filePath)
		if err != nil {
			return ""
		}
	}

	return string(data)
}
