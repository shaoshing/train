package train

type config struct {
	AssetsPath    string
	AssetsURL     string
	PackageAssets bool
}

var Config config = config{
	AssetsPath:    "public",
	AssetsURL:     "/assets",
	PackageAssets: false,
}
