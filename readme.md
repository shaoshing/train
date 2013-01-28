# Train

## Config

### Http

  import "github.com/shaoshing/train"
  http.HandleFunc(train.Config.AssetURL, http.HandlerFunc(train.Handler))


### View

  import . "github.com/shaoshing/train/helpers"
  {{.JavascriptIncludeTag "app"}}
  {{.StylesheetIncludeTag "app"}}

