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

func newProxy(target *url.URL, userHeader string) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		if userHeader != "" {
			res, err := tailscale.WhoIs(req.Context(), req.RemoteAddr)
			if err != nil {
				log.Print("failed to get the owner of the request")
			}
			req.Header.Set(userHeader, res.UserProfile.LoginName)
		}
	}
	return &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			Dial: proxy.FromEnvironment().Dial,
		},
	}
}

func main() {
	t, err := url.Parse(os.Getenv("TPROXY_TARGET_ADDR"))
	if err != nil {
		log.Fatal("failed to parse TPROXY_TARGET_ADDR", err)
	}
	p := newProxy(t, os.Getenv("TPROXY_USER_HEADER"))
	s := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: tailscale.GetCertificate,
		},
		Handler: p,
	}
	log.Print("running tproxy TLS server on :443 ...")
	log.Fatal(s.ListenAndServeTLS("", ""))
}
