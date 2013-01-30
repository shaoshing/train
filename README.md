# Train

## Usages

### Http

```go
  import "github.com/shaoshing/train"
  http.HandleFunc(train.Config.AssetURL, http.HandlerFunc(train.Handler))
```


### Handler


```go
  import "github.com/shaoshing/train"

  type Layout struct{
    Train train.Helpers
  }

  func main() {

    layout := Layout{Train: train.Helpers{}}
    html := `
    {{.Layout.Train.JavascriptTag "app"}}
    {{.Layout.Train.StylesheetTag "app"}}
    `

    tmpl, _ := template.New("").Parse(html)
    tmpl.Execute(os.Stdout, layout)
  }
```

## Production

Install the command line tool to automatically bundle assets into the public folder:

```shell
  go install github.com/shaoshing/train/train
  $GOPATH/bin/train
  ls public/assets
```

Enable bundling assets in Train:

```go
  train.Config.BundleAssets = true
```
