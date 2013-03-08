package train

import (
	"fmt"
)

type config struct {
	AssetsPath   string
	AssetsUrl    string
	BundleAssets bool
	SASS         sassConfig
}

type sassConfig struct {
	DebugInfo   bool
	LineNumbers bool
}

var Config config = config{
	AssetsPath: "assets",
	AssetsUrl:  "/assets/",
}

func init() {
	Config.BundleAssets = IsInProduction()

	if IsInProduction() {
		if err := LoadManifestInfo(); err != nil {
			fmt.Println("== Could not load manifest from public/assets/")
		}
	}
}
