package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/sys/windows/registry"

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

func updateSystemProxySetting() {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	k.SetStringValue("ProxyServer", "127.0.0.1:720")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	updateSystemProxySetting()
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
	log.Fatal(http.ListenAndServe(":720", proxy))
}
