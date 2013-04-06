package train

import (
	"fmt"
	"os"
)

type config struct {
	AssetsPath string
	AssetsUrl  string
	// Show verbose logs. For example, SASS warnings.
	Verbose bool
	// Whether to serve bundled assets in development mode. This option is ignored
	// when in production mode, that is, the ./public/assets folder exists.
	BundleAssets bool
	// When set to DevelopmentMode, assets are read from ./assets
	// When set to ProductionMode, assets are read from ./public/assets
	// It is set to ProductionMode automatically if the ./public/assets exist.
	Mode string
	SASS sassConfig
}

const (
	DevelopmentMode = "development"
	ProductionMode  = "production"
)

type sassConfig struct {
	DebugInfo   bool
	LineNumbers bool
}

var Config = config{
	AssetsPath: "assets",
	AssetsUrl:  "/assets/",
	Mode:       DevelopmentMode,
}

func init() {
	if HasPublicAssets() {
		Config.Mode = ProductionMode
	}

	if IsInProduction() {
		if err := LoadManifestInfo(); err != nil {
			fmt.Println("== Could not load manifest from public/assets/")
		}
	}
}

func IsInProduction() bool {
	return Config.Mode == ProductionMode
}

func HasPublicAssets() bool {
	_, err := os.Stat("public" + Config.AssetsUrl)
	return err == nil
}
