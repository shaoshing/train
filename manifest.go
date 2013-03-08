package train

import (
	"io/ioutil"
	"strings"
)

type FpAssets map[string]string

var manifestInfo FpAssets

const (
	ManifestPath      = "public/assets/manifest.txt"
	ManifestSeparator = "  ->  "
)

func LoadManifestInfo() error {
	manifestInfo = make(FpAssets)

	content, err := ioutil.ReadFile(ManifestPath)
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(content), "\n") {
		info := strings.Split(line, ManifestSeparator)
		if len(info) != 2 {
			continue
		}
		manifestInfo[info[0]] = info[1]
	}
	return nil
}

func WriteToManifest(fpAssets FpAssets) (err error) {
	var content string
	for assetUrl, assetHashedUrl := range fpAssets {
		content += assetUrl + ManifestSeparator + assetHashedUrl + "\n"
	}

	err = ioutil.WriteFile(ManifestPath, []byte(content), 0644)
	return
}
