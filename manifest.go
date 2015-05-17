package train

import (
	"io/ioutil"
	"strings"
)

type FpAssets map[string]string

var ManifestInfo FpAssets

const (
	ManifestSeparator = "  ->  "
)

func ManifestAbsolutePath() string {
	return Config.PublicPath + "/assets/manifest.txt"
}

func LoadManifestInfo() error {
	ManifestInfo = make(FpAssets)

	content, err := ioutil.ReadFile(ManifestAbsolutePath())
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(content), "\n") {
		info := strings.Split(line, ManifestSeparator)
		if len(info) != 2 {
			continue
		}
		ManifestInfo[info[0]] = info[1]
	}
	return nil
}

func WriteToManifest(fpAssets FpAssets) (err error) {
	var content string
	for assetUrl, assetHashedUrl := range fpAssets {
		assetUrl = strings.Replace(assetUrl, Config.PublicPath, "", -1)
		assetHashedUrl = strings.Replace(assetHashedUrl, Config.PublicPath, "", -1)
		content += assetUrl + ManifestSeparator + assetHashedUrl + "\n"
	}

	err = ioutil.WriteFile(ManifestAbsolutePath(), []byte(content), 0644)
	return
}
