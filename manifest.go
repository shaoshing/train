package train

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type FpAssets map[string]string

var manifestInfo = make(FpAssets)

const (
	ManifestPath      = "public/assets/manifest.txt"
	ManifestSeparator = "  ->  "
)

func initManifestInfo() {
	content, err := ioutil.ReadFile(ManifestPath)
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		info := strings.Split(line, ManifestSeparator)
		if len(info) != 2 {
			continue
		}
		manifestInfo[info[0]] = info[1]
	}
}

func WriteToManifest(fpAssets FpAssets) (err error) {
	var content string
	for assetUrl, assetHashedUrl := range fpAssets {
		content += assetUrl + ManifestSeparator + assetHashedUrl + "\n"
	}

	err = ioutil.WriteFile(ManifestPath, []byte(content), 0644)
	return
}
