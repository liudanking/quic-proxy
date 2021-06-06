package main

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/lucas-clemente/quic-go"

	"flag"

	_ "net/http/pprof"

	log "github.com/liudanking/goutil/logutil"
	"github.com/liudanking/quic-proxy/common"
)

func main() {
	var (
		listenAddr string
		cert       string
		key        string
		auth       string
		verbose    bool
		pprofile   bool
	)
	flag.StringVar(&listenAddr, "l", ":443", "listen addr (udp port only)")
	flag.StringVar(&cert, "cert", "", "cert path")
	flag.StringVar(&key, "key", "", "key path")
	flag.StringVar(&auth, "auth", "quic-proxy:Go!", "basic auth, format: username:password")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&pprofile, "p", false, "http pprof")
	flag.Parse()

	log.Info("%v", verbose)
	if cert == "" || key == "" {
		log.Error("cert and key can't by empty")
		return
	}

	if pprofile {
		pprofAddr := "localhost:6060"
		log.Notice("listen pprof:%s", pprofAddr)
		go http.ListenAndServe(pprofAddr, nil)
	}

	parts := strings.Split(auth, ":")
	if len(parts) != 2 {
		log.Error("auth param invalid")
		return
	}
	username, password := parts[0], parts[1]

	listener, err := quic.ListenAddr(listenAddr, generateTLSConfig(cert, key), nil)
	if err != nil {
		log.Error("listen failed:%v", err)
		return
	}
	ql := common.NewQuicListener(listener)

	proxy := goproxy.NewProxyHttpServer()
	ProxyBasicAuth(proxy, func(u, p string) bool {
		return u == username && p == password
	})
	proxy.Verbose = verbose
	server := &http.Server{Addr: listenAddr, Handler: proxy}
	log.Info("start serving %v", listenAddr)
	log.Error("serve error:%v", server.Serve(ql))

}

func generateTLSConfig(certFile, keyFile string) *tls.Config {
	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{common.KQuicProxy},
	}
}
