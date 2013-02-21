# Train

Asset Management Package for web app in Go language. Inspired by [Rails Asset Pipeline](http://guides.rubyonrails.org/asset_pipeline.html).

[![Build Status](https://travis-ci.org/shaoshing/train.png?branch=master)](https://travis-ci.org/shaoshing/train)

## Quick Look

Use Train to manage your asset's dependencies. Enables you to write javascript or stylesheet in the following way:


app.js
```js
//= require jquery

$(function(){
  // Do something cool
})
```

app.css
```css
/*
 *= require reset
 */

body{
}
```

These files will become:

```js
/*! jQuery v@1.8.0 jquery.com | jquery.org/license */
(function(a,b){function G(a){var b=F[a]={};return p.each(a.split(...

$(function(){
  // Do something cool
})
```

```css
html, body, div, span, applet, object, iframe,
h1, h2, h3, h4, h5, h6, p, blockquote, pre,
a, abbr, acronym, address, big, cite, code,
del, dfn, em, img, ins, kbd, q, s, samp,
small, strike, strong, sub, sup, tt, var,
b, u, i, center,
....

body{
}
```

## Usages

### Handler

```go
  import "github.com/shaoshing/train"

  http.HandleFunc(train.Config.AssetURL, http.HandlerFunc(train.Handler))
```


### Template


```go
  import "github.com/shaoshing/train"

  type Layout struct{
    Train train.Helpers // Export Train helpers to templates
  }

  func main() {
    layout := Layout{Train: train.Helpers{}}

    html := `
    {{.Layout.Train.JavascriptTag "app"}}

    {{.Layout.Train.StylesheetTag "app"}}
    `
    tmpl, _ := template.New("").Parse(html)
    tmpl.Execute(os.Stdout, layout)
    //
    // <script src="/assets/javascripts/jquery.js?12345"></script>
    // <script src="/assets/javascripts/app.js?12345"></script>
    //
    // <link rel="stylesheet" href="/assets/stylesheets/reset.css?12345">
    // <link rel="stylesheet" href="/assets/stylesheets/app.css?12345">
  }
```

## Production

Install the command line tool to bundle and compress assets automatically:

```shell
go build -o $GOPATH/bin/train github.com/shaoshing/train/cmd

train
-> clean bundled assets
-> copy assets from assets
-> bundle assets with require directive
-> compress assets

ls public/assets
```

The train tool will bundle your assets into the public/assets folder, with all files expaneded and compressed (by YUI compressor).
You can then use any web servers (nginx, apache, or the Go's file server) to serve these static files.
The template helpers will also stop expanding files if it found the public assets folder. That is, the following code:

```html
{{.Layout.Train.JavascriptTag "app"}}
{{.Layout.Train.StylesheetTag "app"}}
```

Will become:

```html
<script src="/assets/javascripts/app.js?12345"></script>
<link rel="stylesheet" href="/assets/stylesheets/app.css?12345">
```


## License

Train is released under the [MIT License](http://www.opensource.org/licenses/MIT).
