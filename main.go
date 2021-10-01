package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
	"tailscale.com/client/tailscale"
)

func main() {
	t, err := url.Parse(os.Getenv("TPROXY_TARGET_ADDR"))
	if err != nil {
		log.Fatal("failed to parse TPROXY_TARGET_ADDR", err)
	}
	p := httputil.NewSingleHostReverseProxy(t)
	p.Transport = &http.Transport{
		Dial: proxy.FromEnvironment().Dial,
	}
	s := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: tailscale.GetCertificate,
		},
		Handler: p,
	}
	log.Print("running tproxy TLS server on :443 ...")
	log.Fatal(s.ListenAndServeTLS("", ""))
}
