package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	"github.com/nallerooth/http-info/cert"
)

func printLabelValue(indent, label, value string) {
	fmt.Printf("%s%-10s: %s\n", indent, label, value)
}

func printTLSInfo(cs *tls.ConnectionState) {
	indent := "\t"

	fmt.Println("Certificates")
	printLabelValue(indent, "ServerName", cs.ServerName)
	printLabelValue(indent, "Protocol", cs.NegotiatedProtocol)
	fmt.Println()
	for _, c := range cs.PeerCertificates {
		printLabelValue(indent, "Issuer", fmt.Sprintf("%s", c.Issuer))
		printLabelValue(indent, "IsCA", fmt.Sprintf("%t", c.IsCA))
		if c.IsCA == false {
			if len(c.DNSNames) > 1 {
				printLabelValue(indent, "DNSNames", c.DNSNames[0])
				for _, name := range c.DNSNames[1:] {
					printLabelValue(indent, "", name)
				}
			} else {
				printLabelValue(indent, "DNSNames", c.DNSNames[0])
			}
		}
		printLabelValue(indent, "Algorithm", fmt.Sprintf("%s", c.SignatureAlgorithm))
		printLabelValue(indent, "NotBefore", fmt.Sprintf("%s", c.NotBefore))
		notAfterRemaining := fmt.Sprintf("(%d days remaining)", cert.CalcRemainingDays(c.NotAfter))
		printLabelValue(indent, "NotAfter", fmt.Sprintf("%s %s", c.NotAfter, notAfterRemaining))
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

	var timeDNS, timeConn, timeTLS, timeTTFB time.Duration

	fmt.Println("Timings")
	indent := "\t"
	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsDoneInfo = &ddi
			timeDNS = time.Since(dns)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsConnectionState = &cs
			timeTLS = time.Since(tlsHandshake)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			timeConn = time.Since(connect)
		},

		GotFirstResponseByte: func() {
			timeTTFB = time.Since(connect)
		},
	}

	start = time.Now()
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatalln("Request Error:", err)
	}

	if dnsDoneInfo.Err != nil {
		log.Fatalln("DNS Error: ", dnsDoneInfo.Err)
	}

	printLabelValue(indent, "DNS", fmt.Sprintf("%v", timeDNS))
	printLabelValue(indent, "Connect", fmt.Sprintf("%v", timeConn))
	printLabelValue(indent, "TLS", fmt.Sprintf("%v", timeTLS))
	printLabelValue(indent, "TTFB", fmt.Sprintf("%v", timeTTFB))

	printLabelValue(indent, "Total", fmt.Sprintf("%v", time.Since(start)))
	fmt.Println()

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
