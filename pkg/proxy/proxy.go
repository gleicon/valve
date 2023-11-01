package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func echoCertificateDataHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.TLS.PeerCertificates) > 0 {
		fmt.Fprintf(w, "CN=%s\n", r.TLS.PeerCertificates[0].Subject.CommonName)
	} else {
		fmt.Fprintf(w, "No cert")
	}
	//	for _, c := range r.TLS.PeerCertificates {
	//		log.Println(c.Subject)
	//	}
}

func Proxy() {

	remote, err := url.Parse("http://google.com")
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL)
			r.Host = remote.Host
			w.Header().Set("X-Ben", "Rad")
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	// setup handlers and proxy
	http.HandleFunc("/", handler(proxy))
	http.HandleFunc("/echo", echoCertificateDataHandler)
	log.Println("Starting up")

	caCert, err := os.ReadFile("icpcerts/chain.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	//	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	log.Fatal(server.ListenAndServeTLS("mycerts/cert13.pem", "mycerts/privkey13.pem"))

}
