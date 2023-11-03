package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	httplogger "github.com/gleicon/go-httplogger"
)

type ProxyServer struct {
	addr              string
	port              string
	upstream          string
	certFile          string
	keyFile           string
	caCert            string
	echoHandlerActive bool
}

func NewProxyServer(addr, port, upstream, certFile, keyFile, caCert string, echoHandlerActive bool) *ProxyServer {
	ps := ProxyServer{
		addr,
		port,
		upstream,
		certFile,
		keyFile,
		caCert,
		echoHandlerActive,
	}
	return &ps

}

func (ps *ProxyServer) getCNFromRequest(r *http.Request) (string, error) {
	if len(r.TLS.PeerCertificates) > 0 {
		return r.TLS.PeerCertificates[0].Subject.CommonName, nil
	} else {
		return "", errors.New("no certificate")
	}
}

func (ps *ProxyServer) echoCertificateDataHandler(w http.ResponseWriter, r *http.Request) {

	cn, err := ps.getCNFromRequest(r)
	if err != nil {
		fmt.Fprintf(w, "No cert: %v", err)
	}
	fmt.Fprintf(w, "CN=%s\n", cn)

}

func (ps *ProxyServer) setupProxyHandler() func(http.ResponseWriter, *http.Request) {
	remote, err := url.Parse(ps.upstream)
	if err != nil {
		log.Fatal(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Host = remote.Host
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// header overwrite director
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		cn, err := ps.getCNFromRequest(r)

		if err != nil {
			r.Header.Set("X-CERTIFICATE-DETECTED", "off")
		} else {
			r.Header.Set("X-CERTIFICATE-CN", cn)
		}
	}
	return handler(proxy)
}

func (ps *ProxyServer) Serve() {

	serveMux := http.NewServeMux()
	proxyHandler := ps.setupProxyHandler()

	serveMux.HandleFunc("/", proxyHandler)

	// install echo handler if enabled
	if ps.echoHandlerActive {
		serveMux.HandleFunc("/echo", ps.echoCertificateDataHandler)
	}

	log.Println("Reverse Proxy Server Starting up")
	log.Printf("Upstream: %s\n", ps.upstream)

	caCert, err := os.ReadFile(ps.caCert)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()

	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal("Failed loading Cert Pool")
	}

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      ps.addr + ":" + ps.port,
		TLSConfig: tlsConfig,
		Handler:   httplogger.HTTPLogger(serveMux),
	}

	log.Fatal(server.ListenAndServeTLS(ps.certFile, ps.keyFile))

}
