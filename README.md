# http-info
A simple HTTP information tool for the command line.

**Status:** In development.

## Example
```
$ http-info https://github.com

DNS
	Resolved IPs:
	          : 140.82.121.4

Timings
	DNS       : 6.447103ms
	Connect   : 37.815259ms
	TLS       : 64.571003ms
	TTFB      : 151.960028ms
	Total     : 159.172329ms

Transfer
	Status    : 200 OK
	Bytes     : 239269
	Compressed: true

Certificates
	ServerName: github.com
	Protocol  : h2

	Issuer    : CN=DigiCert High Assurance TLS Hybrid ECC SHA256 2020 CA1,O=DigiCert\, Inc.,C=US
	IsCA      : false
	DNSNames  : github.com
	          : www.github.com
	Algorithm : ECDSA-SHA256
	NotBefore : 2021-03-25 00:00:00 +0000 UTC
	NotAfter  : 2022-03-30 23:59:59 +0000 UTC [ 182 days remaining ]

	Issuer    : CN=DigiCert High Assurance EV Root CA,OU=www.digicert.com,O=DigiCert Inc,C=US
	IsCA      : true
	Algorithm : SHA256-RSA
	NotBefore : 2020-12-17 00:00:00 +0000 UTC
	NotAfter  : 2030-12-16 23:59:59 +0000 UTC [ 3365 days remaining ]

Done
```

## Building/Installing
Developed in Go `1.16`.

```sh
go build
```
or
```sh
go install
```
