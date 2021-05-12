package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func printLabelValue(indent, label, value string) {
	fmt.Printf("%s%-10s: %s\n", indent, label, value)
}

func printTLSInfo(cs *tls.ConnectionState) {
	indent := "\t"
	//fmt.Printf("%sTLS: %+v\n", indent, cs)

	fmt.Println("Certificates")
	printLabelValue(indent, "ServerName", cs.ServerName)
	fmt.Println()
	for _, c := range cs.PeerCertificates {
		printLabelValue(indent, "Issuer", fmt.Sprintf("%s", c.Issuer))
		printLabelValue(indent, "IsCA", fmt.Sprintf("%t", c.IsCA))
		if c.IsCA == false {
			printLabelValue(indent, "DNSNames", fmt.Sprintf("%s", c.DNSNames))
		}
		//printLabelValue(indent, "SignatureAlgorithm", fmt.Sprintf("%s", c.SignatureAlgorithm))
		printLabelValue(indent, "NotBefore", fmt.Sprintf("%s", c.NotBefore))
		printLabelValue(indent, "NotAfter", fmt.Sprintf("%s", c.NotAfter))
		fmt.Println()
	}
}

func printDNSInfo(ddi *httptrace.DNSDoneInfo) {
	indent := "\t"

	fmt.Println("DNS")
	printLabelValue(indent, "Resolved IPs", "")
	for _, ip := range ddi.Addrs {
		printLabelValue(indent, "", ip.String())
	}
	fmt.Println()
}

func timeGet(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns, tlsHandshake time.Time
	var tlsConnectionState *tls.ConnectionState
	var dnsDoneInfo *httptrace.DNSDoneInfo

	indent := ""

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsDoneInfo = &ddi
			fmt.Printf("%sDNS Done: %v\n", indent, time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsConnectionState = &cs
			fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("Connect time: %v\n", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			fmt.Printf("Time from start to first byte: %v\n", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))

	printDNSInfo(dnsDoneInfo)
	printTLSInfo(tlsConnectionState)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Missing arguments")
	}

	url := args[0]
	timeGet(url)

	fmt.Println("Done")
}
