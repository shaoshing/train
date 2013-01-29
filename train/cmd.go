package main

import (
	"github.com/shaoshing/train"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	removeAssets()
	copyAssets()
	bundleAssets()
}

func removeAssets() {
	err := exec.Command("rm", "-rf", "public"+train.Config.AssetsUrl).Run()
	if err != nil {
		panic(err)
	}
}

func copyAssets() {
	err := exec.Command("cp", "-rf", train.Config.AssetsPath, "public"+train.Config.AssetsUrl).Run()
	if err != nil {
		panic(err)
	}
}

func bundleAssets() {
	train.Config.BundleAssets = true
	publicAssetPath := "public" + train.Config.AssetsUrl
	filepath.Walk(publicAssetPath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileExt := path.Ext(filePath)
		if fileExt != ".js" && fileExt != ".css" {
			return nil
		}

		assetUrl := strings.Replace(filePath, publicAssetPath, train.Config.AssetsUrl, 1)
		b_content, err := ioutil.ReadFile(filePath)
		content := string(b_content)
		header := train.FindDirectivesHeader(&content, fileExt)
		if len(header) != 0 {
			content := train.ReadAsset(assetUrl)
			ioutil.WriteFile(filePath, []byte(content), os.ModeDevice)
		}
		return nil
	})
}
