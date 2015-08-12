package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

func isBaidu() goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		log.Print("Host:" + req.URL.Host)
		log.Print("path:" + req.URL.Path)
		log.Print("schema:" + req.URL.Scheme)
		return strings.Contains(req.URL.Host, "oms.com")
	}
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest(isBaidu()).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		log.Print("url:" + r.URL.String())
		res, err := http.Get("http://127.0.0.1:8080" + r.URL.Path)
		if err != nil {
			fmt.Println(err.Error())
		}
		return r, res
	})
	log.Fatal(http.ListenAndServe(":8081", proxy))
}
