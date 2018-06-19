Introduction
============

This is an extremely minimal router, for usage in your own application. Supports only 2 things:
* routing based on request-method and path
* parameters in URLs

I build it for a json-only REST api with a few routes.

Usage
=====

This intentionally does NOT implement `ServeHTTP`, as you should implement this
yourself, to do any wrapping and processing of request or response.

Example
=======

```go
package main

import (
	"encoding/json"
	"github.com/parse-nl/MinimalRouter"
	"net/http"
)

type appController struct{}
type appRequest    struct{ Params map[string]string }
type appResponse   struct{ Result string }

var router *minimalrouter.Router

func init() {
	router = minimalrouter.New()

	router.Add("GET",  "/ping",             handleGetPing)
	router.Add("POST", "/hello-from/:name", handlePostHello)
}

func main() {
	http.ListenAndServe(":8088", &appController{})
}

// implement the standard ServeHTTP as required by golang
// handles json encoding and decoding, could implement authentication as well
func (c *appController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var result *appResponse

	fn, params := router.Match(r.Method, r.URL.Path)
	request := &appRequest{params}

	if fn == nil {
		w.WriteHeader(http.StatusNotFound)
		result = &appResponse{"error"}
	} else {
		w.WriteHeader(http.StatusOK)

		// convert the generic interface{} to the function-signature we expect
		result = fn.(func(*appRequest) *appResponse)(request)
	}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(result)
}

func handleGetPing(r *appRequest) *appResponse {
	return &appResponse{"pong"}
}

func handlePostHello(r *appRequest) *appResponse {
	return &appResponse{"hey " + r.Params["name"]}
}

```