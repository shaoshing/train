# Train [![Build Status](https://travis-ci.org/shaoshing/train.png?branch=master)](https://travis-ci.org/shaoshing/train)

Asset Management Package for web app in Go language. Inspired by [Rails Asset Pipeline](http://guides.rubyonrails.org/asset_pipeline.html).

## Concepts

Most of the ideas in train are borrowed from [Rails' Assets Pipeline](http://guides.rubyonrails.org/asset_pipeline.html) which is built on top of [Sprockets](https://github.com/sstephenson/sprockets). Train will compile, compress and fingerprint your assets. However, the actual implementation of train is mainly different from it, due to the reason of simplification. Unlike Pipeline, train will not try to grab so much control of assets path.

Like Pipeline, there is also the existence of manifest file, which is very important and very convenient to manage your assets' dependency.

### Manifest File

You need add at least one template tag in your html page when you want to ask train to take care of your assets management, like this:

```go
{{javascript_tag "app"}}
```

And by default, in your `assets/javascripts/` folder, there must be a `app.js` or `app.coffee` file. Now the file will be treated as a manifest file.

What can we do in a manifest file? Add `require` directives(See more at Quick Look section).

### App Layout

```
|- app.go
|
|- assets
	|- javascripts
	|- stylesheets
	|- images
|- public
	|- assets
		|- javascripts
		|- stylesheets
		|- images
```

By default, you should put all your Javascript assets inside the `assets/javascripts` folder, css inside `assets/stylesheets`. Even though there are not limitations about your image assets, but as most good fellows did, `assets/images` is the right way to go.

`public/assets` is an important trick `train` uses. Each time when you boot train, it will detect whether the `public/assets` exits, if it does, train will serve your assets from that folder instead of your `assets` folder. Since then, when your app is still under developing, remember to remove your `public/assets` folder and don't put any of your files inside this folder in case of misguiding `train` or some unnecessary loss of time and codes.

## Quick Look

Use Train to manage your asset's dependencies. Enables you to write javascript or stylesheet in the following way:

### Javascript

assets/javascripts/base.js

```js
appName = "...";
```

assets/javascripts/app.js

```js
//= require javascripts/base

$(function(){
  // Do something cool
})
```

GET /assets/javascripts/app.js

```js
appName = "...";

$(function(){
  // Do something cool
})
```

### Stylesheet

assets/stylesheets/base.css

```css
h1, h2{ padding:0; }
```

assets/stylesheets/app.css

```css
/*
 *= require stylesheets/base
 */

body{...}
```

GET /assets/stylesheets/app.css

```css
h1, h2{ padding:0; }

body{...}
```

### CoffeeScript

assets/javascripts/app.coffee

```coffee
alert "Hello CoffeeScript!"
```

GET /assets/javascripts/app.js

```coffee
alert("Hello CoffeeScript!");
```

### SASS

assets/stylesheets/app.sass

```sass
body
  color: red
```

GET /assets/stylesheets/app.js

```sass
body{
  color: red; }
```

## Usage

### Typical Usage

By default, train will register '/assets/' http url pattern, if you want to use another url prefix, please re-config it before calling `Run`.

```go
  package main

  import (
    "fmt"
    "github.com/shaoshing/train"
    "net/http"
  )

  func main() {
      // By passing the nil param, train will add http handler to the http.DefaultServeMux to handle all asset requests starting with "/assets/". However, if you are using custom ServeMux, you can simply pass your ServeMux into it.
      train.ConfigureHttpHandler(nil)
      train.ConfigureHttpHandler(nil)

      fmt.Println("Listening to localhost:8000")
      http.ListenAndServe(":8000", nil)
  }
```

### Default Configurations

Train package has exported `Config` as the configurational interface. currently there are two options:

```go
{
	AssetsPath: "assets",
	AssetsUrl:  "/assets/",
	// Whether to serve bundled assets in development mode. This option is ignored
	// when in production mode, that is, the ./public/assets folder exists.
	BundleAssets bool
	// When set to DevelopmentMode, assets are read from ./assets
	// When set to ProductionMode, assets are read from ./public/assets
	// It is set to ProductionMode automatically if the ./public/assets exist.
	Mode string
	SASS sassConfig
}
```

`AssetsPath` is used to specify the the folder where you put all your assets, defaulted to "assets".

`AssetsUrl` is used to specify the url pattern that train will use to register on http.ServeMux, defaulted to "/assets/"

If you want to change the default configuration, make sure you change it before `train.ConfigureHttpHandler`.

### Re-config Assets URL

```go
package main

import (
  "fmt"
  "github.com/shaoshing/train"
  "net/http"
)

func main() {
    train.Config.AssetsUrl = "/custom/path"
    train.ConfigHttpHandler(nil)

    fmt.Println("Listening to localhost:8000")
    http.ListenAndServe(":8000", nil)
}
```

### Template Helpers

In order to use `javascript_tag` and `stylesheet_tag`, you need add train's template helpers. Like below:

```go
  import "github.com/shaoshing/train"

  func main() {
    tmpl := template.New("index")
    tmpl.Funcs(train.HelperFuncs)
    tmpl.Parse(`
    {{define "index"}}
      {{javascript_tag "app"}}

      {{stylesheet_tag "app"}}
    {{end}}
    `)

    tmpl.Execute(os.Stdout, "index", nil)
    //
    // <script src="/assets/javascripts/base.js?12345"></script>
    // <script src="/assets/javascripts/app.js?12345"></script>
    //
    // <link rel="stylesheet" href="/assets/stylesheets/base.css?12345">
    // <link rel="stylesheet" href="/assets/stylesheets/app.css?12345">
  }
```

## Production

Install the command line tool to bundle and compress assets automatically:

```sh
# Install the train command. Run "train help" to see help info.
# Be sure to add $GOPATH/bin to you $PATH env. Otherwise you will have to run $GOPATH/bin/train
go build -o $GOPATH/bin/train github.com/shaoshing/train/cmd

train
-> clean bundled assets
-> copy assets from assets
-> bundle assets with require directive
-> compress assets

ls public/assets
```

The train command will bundle your assets into the public/assets folder, with all files expaneded and compressed (by YUI compressor).
You can then use any web servers (nginx, apache, or the Go's file server) to serve these static files.
The template helpers will also stop expanding files if it found the public assets folder. That is, the following code:

```go
{{javascript_tag "app"}}
{{stylesheet_tag "app"}}
```

Will become:

	<script src="/assets/javascripts/app-f72c58d3009ff4412f20393e0447674c.js"></script>
	<link rel="stylesheet" href="/assets/stylesheets/app-df05cdccb878a7efa14a98ea2e34e894.css">

## License

Train is released under the [MIT License](http://www.opensource.org/licenses/MIT).
