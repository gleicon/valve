package main

import (
	"flag"

	"github.com/gleicon/valve-go/pkg/proxy"
)

func main() {
	var port = flag.String("p", "8443", "Port to bind")
	var upstream = flag.String("u", "https://google.com", "upstream server")
	var certFile = flag.String("c", "mycerts/cert.pem", "Certfile (.pem)")
	var keyFile = flag.String("k", "mycerts/key.pem", "Private key")
	var caCert = flag.String("a", "icpcerts/chain.pem", "CACerts for ICP Brasil (or another CA Cert you want to ask the client)")
	flag.Parse()

	ps := proxy.NewProxyServer(*port, *upstream, *certFile, *keyFile, *caCert, true)

	ps.Serve()
}
