package main

import (
	"flag"
	"net/http"
	"net/url"

	"github.com/elazarl/goproxy"
	log "github.com/liudanking/goutil/logutil"
	"github.com/liudanking/quic-proxy/common"
)

func main() {
	log.Debug("client")

	var (
		listenAddr     string
		proxyUrl       string
		skipCertVerify bool
		verbose        bool
	)

	flag.StringVar(&listenAddr, "l", ":18080", "listenAddr")
	flag.StringVar(&proxyUrl, "proxy", "", "upstream proxy url")
	flag.BoolVar(&skipCertVerify, "k", false, "skip Cert Verify")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	Url, err := url.Parse(proxyUrl)
	if err != nil {
		log.Error("proxyUrl:%s invalid", proxyUrl)
		return
	}
	if Url.Scheme == "https" {
		log.Error("quic-proxy only support http proxy")
		return
	}

	proxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}

	dialer := common.NewQuicDialer(skipCertVerify)
	proxy.Tr.Dial = dialer.Dial

	proxy.ConnectDial = proxy.NewConnectDialToProxy(proxyUrl)

	log.Info("start serving %s", listenAddr)
	log.Error("%v", http.ListenAndServe(listenAddr, proxy))
}
