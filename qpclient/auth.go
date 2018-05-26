package main

import (
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/liudanking/quic-proxy/common"
)

func SetAuthForBasicRequest(username, password string) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		common.SetBasicAuth(username, password, req)
		return req, nil
	})
}

func SetAuthForBasicConnectRequest(username, password string) func(req *http.Request) {
	return func(req *http.Request) {
		common.SetBasicAuth(username, password, req)
	}
}
