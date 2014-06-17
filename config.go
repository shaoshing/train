package train

import (
	"fmt"
	"os"
)

type config struct {
	// The folder that contains JavaScript and StyleSheet files. It will use "assets" by default.
	AssetsPath string
	AssetsUrl  string
	// Show verbose logs. For example, SASS warnings.
	Verbose bool
	// Whether to serve bundled assets in development mode. This option is ignored
	// when in production mode, that is, the ./public/assets folder exists.
	BundleAssets bool
	// When set to DEVELOPMENT_MODE, assets are read from ./assets
	// When set to PRODUCTION_MODE, assets are read from ./public/assets
	// It is set to PRODUCTION_MODE automatically if the ./public/assets exist.
	Mode string
	SASS sassConfig
}

const (
	DEVELOPMENT_MODE = "development"
	PRODUCTION_MODE  = "production"
	VERSION          = "0.1"
)

type sassConfig struct {
	DebugInfo   bool
	LineNumbers bool
}

var Config = config{
	AssetsPath: "assets",
	AssetsUrl:  "/assets/",
	Mode:       DEVELOPMENT_MODE,
}

func init() {
	if HasPublicAssets() {
		Config.Mode = PRODUCTION_MODE
	}

	if IsInProduction() {
		if err := LoadManifestInfo(); err != nil {
			fmt.Println("== Could not load manifest from public/assets/")
		}
	}
}

func IsInProduction() bool {
	return Config.Mode == PRODUCTION_MODE
}

func HasPublicAssets() bool {
	_, err := os.Stat("public" + Config.AssetsUrl)
	return err == nil
}
