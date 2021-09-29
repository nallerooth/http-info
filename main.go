package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"strings"
	"time"

	"github.com/nallerooth/http-info/cert"
	"github.com/nallerooth/http-info/colors"
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
		notAfterRemaining := fmt.Sprintf("[ %s ]", cert.CalcRemainingDaysColor(c.NotAfter))
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
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatalln("Request Error:", err)
	}
	defer res.Body.Close()

	if dnsDoneInfo.Err != nil {
		log.Fatalln("DNS Error: ", dnsDoneInfo.Err)
	}

	printDNSInfo(dnsDoneInfo)
	fmt.Println("Timings")

	printLabelValue(indent, "DNS", fmt.Sprintf("%v", timeDNS))
	printLabelValue(indent, "Connect", fmt.Sprintf("%v", timeConn))
	printLabelValue(indent, "TLS", fmt.Sprintf("%v", timeTLS))
	printLabelValue(indent, "TTFB", fmt.Sprintf("%v", timeTTFB))

	printLabelValue(indent, "Total", fmt.Sprintf("%v", time.Since(start)))
	fmt.Println()

	fmt.Println("Transfer")
	// Write response to /dev/null and count number of bytes written
	dn, err := os.OpenFile(os.DevNull, os.O_WRONLY, os.FileMode(fs.ModeAppend))
	if err != nil {
		log.Fatalf("Error opening DevNull (%s): %s", os.DevNull, err)
	}
	defer dn.Close()
	numBytes, err := io.Copy(dn, res.Body)

	colorFn := colors.White
	switch {
	case res.StatusCode >= 400:
		colorFn = colors.Red
	case res.StatusCode >= 300:
		colorFn = colors.Yellow
	case res.StatusCode >= 200:
		colorFn = colors.Green
	case res.StatusCode >= 100:
		colorFn = colors.Purple
	}

	printLabelValue(indent, "Status", colorFn(res.Status))
	redirects := map[int]bool{
		301: true,
		302: true,
		303: true,
		307: true,
	}
	if redirects[res.StatusCode] && res.Header.Get("Location") != "" {
		printLabelValue(indent, "Redirect", res.Header.Get("Location"))
	}
	printLabelValue(indent, "Bytes", fmt.Sprintf("%d", numBytes))
	printLabelValue(indent, "Compressed", fmt.Sprintf("%t", res.Uncompressed))
	if len(res.TransferEncoding) > 0 {
		printLabelValue(indent, "Encoding", strings.Join(res.TransferEncoding, ", "))
	}
	fmt.Println()

	if tlsConnectionState != nil {
		printTLSInfo(tlsConnectionState)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic("Missing arguments")
	}

	url := args[0]
	fmt.Println()
	timeGet(url)
	fmt.Println("Done")
}
