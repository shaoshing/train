package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/shaoshing/train"
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	if !prepareEnv() {
		return
	}
	removeAssets()
	copyAssets()
	bundleAssets()
	compressAssets()
	fingerPrintAssets()
	train.Stop()
}

func prepareEnv() bool {
	public, err := os.Stat("public")

	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir("public", os.FileMode(0777))
		if err != nil {
			panic(err)
		}
	} else if !public.IsDir() {
		fmt.Println("Can't create public directory automatically because a file with the same name already exists.\nPlease consider renaming your file or moving it to another folder.")
		return false
	}

	return true
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

var mapCompiledExt = map[string]string{
	".sass":   ".css",
	".scss":   ".css",
	".coffee": ".js",
}

func bundleAssets() {
	fmt.Println("-> bundle and compile assets")

	train.Config.BundleAssets = true
	publicAssetPath := "public" + train.Config.AssetsUrl

	filepath.Walk(publicAssetPath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		assetUrl := strings.Replace(filePath, publicAssetPath, train.Config.AssetsUrl, 1)
		fileExt := path.Ext(filePath)
		switch fileExt {
		case ".js", ".css":
			if hasRequireDirectives(filePath) {
				content, err := train.ReadAsset(assetUrl)
				if err != nil {
					removeAssets()
					panic(err)
				}
				ioutil.WriteFile(filePath, []byte(content), os.ModeDevice)
			}
		case ".sass", ".scss", ".coffee":
			if path.Base(filePath)[0] == '_' {
				return nil
			}

			content, err := train.ReadAsset(assetUrl)
			if err != nil {
				removeAssets()
				fmt.Println("Error when reading asset: ", assetUrl)
				panic(err)
			}
			compiledPath := strings.Replace(filePath, fileExt, mapCompiledExt[fileExt], 1)
			os.Create(compiledPath)
			ioutil.WriteFile(compiledPath, []byte(content), os.ModeDevice)
		default:
			return nil
		}
		return nil
	})
}

func hasRequireDirectives(filePath string) bool {
	b_content, _ := ioutil.ReadFile(filePath)
	content := string(b_content)
	fileExt := path.Ext(filePath)
	header := train.FindDirectivesHeader(&content, fileExt)
	return len(header) != 0
}

func compressAssets() {
	fmt.Println("-> compress assets")

	jsFiles, cssFiles := getCompiledAssets(regexp.MustCompile(`(min\.\w+$)|\/min\/`))

	compress(jsFiles, ".js$:.js")
	compress(cssFiles, ".css$:.css")
}

func compress(files []string, option string) {
	_, err := exec.LookPath("java")
	if err != nil {
		fmt.Println("You don't have Java installed.")
		return
	}

	_, filename, _, _ := runtime.Caller(1)
	pkgPath := path.Dir(filename)
	yuicompressor := pkgPath + "/yuicompressor-2.4.7.jar"
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

func getCompiledAssets(filter *regexp.Regexp) (jsFiles []string, cssFiles []string) {
	publicAssetPath := "public" + train.Config.AssetsUrl
	filepath.Walk(publicAssetPath, func(filePath string, info os.FileInfo, err error) error {
		fileExt := path.Ext(filePath)
		if filter != nil && filter.Match([]byte(filePath)) {
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

	return
}

func fingerPrintAssets() {
	fmt.Println("-> Fingerprinting Assets")

	assets, cssFiles := getCompiledAssets(nil)
	for _, file := range cssFiles {
		assets = append(assets, file)
	}

	fpAssets := train.FpAssets{}
	for _, asset := range assets {
		assetContent, err := ioutil.ReadFile(asset)
		if err != nil {
			fmt.Printf("Fingerprint Error: %s\n", err)
			return
		}

		h := md5.New()
		io.WriteString(h, string(assetContent))
		fpStr := string(h.Sum(nil))

		dir, file := filepath.Split(asset)
		ext := filepath.Ext(file)
		filename := filepath.Base(file)
		filename = filename[0:strings.LastIndex(filename, ext)]

		fpAsset := fmt.Sprintf("%s%s-%x%s", dir, filename, fpStr, ext)

		err = ioutil.WriteFile(fpAsset, assetContent, 0644)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		fpAssets[asset[6:]] = fpAsset[6:]
	}

	d, err := goyaml.Marshal(&fpAssets)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("public/assets/manifest.yml", d, 0644)
	if err != nil {
		panic(err)
	}
}
