package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/proxy"
	"tailscale.com/client/tailscale"
)

type pathConfig struct {
	target *url.URL
	path   string
}

func newProxy(pathConfigs []pathConfig, userHeader string) *httputil.ReverseProxy {
	for _, config := range pathConfigs {
		log.Printf("%+v", config)
	}
	director := func(req *http.Request) {
		for _, config := range pathConfigs {
			if strings.HasPrefix(req.URL.Path, config.path) {
				req.URL.Scheme = config.target.Scheme
				req.URL.Host = config.target.Host
				break
			}
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

func parsePathConfigs() []pathConfig {
	if targetAddr, ok := os.LookupEnv("TPROXY_TARGET_ADDR"); ok {
		target, err := url.Parse(targetAddr)
		if err != nil {
			log.Fatal("failed to parse TPROXY_TARGET_ADDR", err)
		}
		return []pathConfig{{target, ""}}
	}
	if config, ok := os.LookupEnv("TPROXY_TARGET_PATH_CONFIGS"); ok {
		configStrings := strings.Split(config, ",")
		pathConfigs := make([]pathConfig, len(configStrings))
		for i, configString := range configStrings {
			subStrings := strings.SplitN(configString, "#", 2)
			path := ""
			if len(subStrings) == 2 {
				path = subStrings[0]
			}
			target, err := url.Parse(subStrings[len(subStrings)-1])
			if err != nil {
				log.Fatal("failed to parse TPROXY_TARGET_PATH_CONFIGS", err)
			}
			pathConfigs[i] = pathConfig{target, path}
		}
		return pathConfigs
	}
	log.Fatal("no target address specified, use the TPROXY_TARGET_ADDR or TPROXY_TARGET_PATH_CONFIGS env variables")
	return nil
}

func main() {
	p := newProxy(parsePathConfigs(), os.Getenv("TPROXY_USER_HEADER"))
	s := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: tailscale.GetCertificate,
		},
		Handler: p,
	}
	log.Print("running tproxy TLS server on :443 ...")
	log.Fatal(s.ListenAndServeTLS("", ""))
}
