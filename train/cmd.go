package main

import (
	"bytes"
	"fmt"
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
	compressAssets()
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
			content, err := train.ReadAsset(assetUrl)
			if err != nil {
				removeAssets()
				panic(err)
			}
			ioutil.WriteFile(filePath, []byte(content), os.ModeDevice)
		}
		return nil
	})
}

func compressAssets() {
	var jsFiles, cssFiles []string
	publicAssetPath := "public" + train.Config.AssetsUrl
	filepath.Walk(publicAssetPath, func(filePath string, info os.FileInfo, err error) error {
		fileExt := path.Ext(filePath)
		switch fileExt {
		// Skip minified files
		case ".js":
			jsFiles = append(jsFiles, filePath)
		case ".css":
			cssFiles = append(cssFiles, filePath)
		}
		return nil
	})

	compress(jsFiles, ".js$:.js")
	compress(cssFiles, ".css$:.css")
}

func compress(files []string, option string) {
	yuicompressor := os.Getenv("GOPATH") + "/src/github.com/shaoshing/train/train/yuicompressor-2.4.7.jar"
	cmd := exec.Command("sh", "-c", "java -jar "+yuicompressor+" -o '"+option+"' "+strings.Join(files, " "))
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("YUI Compressor error:", out.String())
		panic(err)
	}
}
