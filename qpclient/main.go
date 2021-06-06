package main

import (
	"flag"
	"net/http"
	"net/url"
	"strings"

	_ "net/http/pprof"

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
		auth           string
		verbose        bool
		pprofile       bool
	)

	flag.StringVar(&listenAddr, "l", ":18080", "listenAddr")
	flag.StringVar(&proxyUrl, "proxy", "", "upstream proxy url")
	flag.BoolVar(&skipCertVerify, "k", false, "skip Cert Verify")
	flag.StringVar(&auth, "auth", "quic-proxy:Go!", "basic auth, format: username:password")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&pprofile, "p", false, "http pprof")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	if pprofile {
		pprofAddr := "localhost:6061"
		log.Notice("listen pprof:%s", pprofAddr)
		go http.ListenAndServe(pprofAddr, nil)
	}

	Url, err := url.Parse(proxyUrl)
	if err != nil {
		log.Error("proxyUrl:%s invalid", proxyUrl)
		return
	}
	if Url.Scheme == "https" {
		log.Error("quic-proxy only support http proxy")
		return
	}

	parts := strings.Split(auth, ":")
	if len(parts) != 2 {
		log.Error("auth param invalid")
		return
	}
	username, password := parts[0], parts[1]

	proxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl)
	}

	dialer := common.NewQuicDialer(skipCertVerify)
	proxy.Tr.Dial = dialer.Dial

	// proxy.ConnectDial = proxy.NewConnectDialToProxy(proxyUrl)
	proxy.ConnectDial = proxy.NewConnectDialToProxyWithHandler(proxyUrl,
		SetAuthForBasicConnectRequest(username, password))

	// set basic auth
	proxy.OnRequest().Do(SetAuthForBasicRequest(username, password))

	log.Info("start serving %s", listenAddr)
	log.Error("%v", http.ListenAndServe(listenAddr, proxy))
}
