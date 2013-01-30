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
	"regexp"
	"strings"
)

func main() {
	removeAssets()
	copyAssets()
	bundleAssets()
	compressAssets()
}

func removeAssets() {
	fmt.Println("-> clean bundled assets")
	err := exec.Command("rm", "-rf", "public"+train.Config.AssetsUrl).Run()
	if err != nil {
		panic(err)
	}
}

func copyAssets() {
	fmt.Println("-> copy assets from", train.Config.AssetsPath)
	err := exec.Command("cp", "-rf", train.Config.AssetsPath, "public"+train.Config.AssetsUrl).Run()
	if err != nil {
		panic(err)
	}
}

func bundleAssets() {
	fmt.Println("-> bundle assets with require directive")
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

var minifiedFiles = regexp.MustCompile(`(min\.\w+$)|\/min\/`)

func compressAssets() {
	fmt.Println("-> compress assets")
	var jsFiles, cssFiles []string
	publicAssetPath := "public" + train.Config.AssetsUrl
	filepath.Walk(publicAssetPath, func(filePath string, info os.FileInfo, err error) error {
		fileExt := path.Ext(filePath)
		if minifiedFiles.Match([]byte(filePath)) {
			return nil
		}
		switch fileExt {
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
	_, err := exec.LookPath("java")
	if err != nil {
		fmt.Println("You don't have Java installed.")
		return
	}

	yuicompressor := os.Getenv("GOPATH") + "/src/github.com/shaoshing/train/train/yuicompressor-2.4.7.jar"
	cmd := exec.Command("sh", "-c", "java -jar "+yuicompressor+" -o '"+option+"' "+strings.Join(files, " "))
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out

	fmt.Println(files)

	err = cmd.Run()
	if err != nil {
		fmt.Println("YUI Compressor error:", out.String())
		panic(err)
	}
}
