# 2013-3-8

* interpreter memory leaking problem solved
* fingerprinting on production mode
* some internal refactoring

# 2013-3-4

* HasPublicAssets is renamed into IsInProduction for codes readablity.
* Using new package using interface:

```go
train.Run()

// instead of
http.Handle(train.Config.AssetsUrl, http.HandlerFunc(train.Handler))
```

* Log assets request by default.

