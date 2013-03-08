package train

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

type FpAssets map[string]string

var manifestInfo FpAssets

func initManifestInfo() {
	dat, err := ioutil.ReadFile("public/assets/manifest.yml")
	if err != nil {
		panic(err)
	}

	goyaml.Unmarshal(dat, &manifestInfo)
}
