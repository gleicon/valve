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
)

type ProxyServer struct {
	port              string
	upstream          string
	certFile          string
	keyFile           string
	caCert            string
	echoHandlerActive bool
}

func NewProxyServer(port string, upstream, certFile, keyFile, caCert string, echoHandlerActive bool) *ProxyServer {
	ps := ProxyServer{
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
			log.Println(r.URL)
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

	proxyHandler := ps.setupProxyHandler()
	http.HandleFunc("/", proxyHandler)

	// optional echo handler
	if ps.echoHandlerActive {
		http.HandleFunc("/echo", ps.echoCertificateDataHandler)
	}

	log.Println("Proxy Server Starting up !")

	caCert, err := os.ReadFile(ps.caCert)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      ":" + ps.port,
		TLSConfig: tlsConfig,
	}

	log.Fatal(server.ListenAndServeTLS(ps.certFile, ps.keyFile))

}
