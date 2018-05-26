package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	log "github.com/liudanking/goutil/logutil"
	"github.com/liudanking/quic-proxy/common"
)

var unauthorizedMsg = []byte("404 not found")

func BasicUnauthorized(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode:    404, // for purpose of avoiding proxy detection
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       req,
		Header:        http.Header{},
		Body:          ioutil.NopCloser(bytes.NewBuffer(unauthorizedMsg)),
		ContentLength: int64(len(unauthorizedMsg)),
	}
}

func auth(req *http.Request, f func(user, passwd string) bool) bool {
	authheader := strings.SplitN(req.Header.Get(common.ProxyAuthHeader), " ", 2)
	req.Header.Del(common.ProxyAuthHeader)
	if len(authheader) != 2 || authheader[0] != "Basic" {
		return false
	}
	userpassraw, err := base64.StdEncoding.DecodeString(authheader[1])
	if err != nil {
		return false
	}
	userpass := strings.SplitN(string(userpassraw), ":", 2)
	if len(userpass) != 2 {
		return false
	}
	return f(userpass[0], userpass[1])
}

// Basic returns a basic HTTP authentication handler for requests
//
// You probably want to use auth.ProxyBasic(proxy) to enable authentication for all proxy activities
func Basic(f func(user, passwd string) bool) goproxy.ReqHandler {
	return goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !auth(req, f) {
			log.Warning("basic auth verify for normal request failed")
			return nil, BasicUnauthorized(req)
		}
		req.Header.Del(common.ProxyAuthHeader)
		return req, nil
	})
}

// BasicConnect returns a basic HTTP authentication handler for CONNECT requests
//
// You probably want to use auth.ProxyBasic(proxy) to enable authentication for all proxy activities
func BasicConnect(f func(user, passwd string) bool) goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if !auth(ctx.Req, f) {
			log.Warning("basic auth verify for connect request failed")
			ctx.Resp = BasicUnauthorized(ctx.Req)
			return goproxy.RejectConnect, host
		}
		ctx.Req.Header.Del(common.ProxyAuthHeader)
		return goproxy.OkConnect, host
	})
}

// ProxyBasic will force HTTP authentication before any request to the proxy is processed
func ProxyBasicAuth(proxy *goproxy.ProxyHttpServer, f func(user, passwd string) bool) {
	proxy.OnRequest().Do(Basic(f))
	proxy.OnRequest().HandleConnect(BasicConnect(f))
}
