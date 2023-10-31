package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
	log.Println(r.TLS.PeerCertificates)
	for _, c := range r.TLS.PeerCertificates {
		log.Println(*&c.Subject)
	}

	log.Println(r.TLS)

}

func Proxy() {
	http.HandleFunc("/hello", helloHandler)
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
