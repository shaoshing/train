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
    {{.Layout.Train.JavascriptIncludeTag "app"}}
    {{.Layout.Train.StylesheetIncludeTag "app"}}
    `

    tmpl, _ := template.New("").Parse(html)
    tmpl.Execute(os.Stdout, layout)
  }
```
