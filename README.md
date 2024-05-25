# No2_GSLC2_NetSysProg
client.go and server.go 
Go program retrieve web resources using the default HTTP client. Include the following features:
1. retrieve web resources using default HTTP client
2.  Implement functionality to close the response body after retrieving resources to prevent resource leaks.
3.  Add timeouts and cancellation mechanisms to the HTTP client to handle cases of slow or unresponsive servers.
4.  Disable persistent TCP connections in the HTTP client configuration.
5.  Support posting data (JSON) over HTTP

Generate.go
is program to Generate key and sertificate to localhost web
*how do i create https connection in localhost : generate key and sertificate in generate.go using following command go run generate.go -host=localhost -> Import cert.pem to browser(chrome) open chrome://settings/privacy -> security -> manage certificate -> import that cert.pem in Trusted Root Certification Authorities 

viewTLS_versionCiphersuiteNameAndIssuer.go 
Is a program to do TCP dial to the web that uses HTTPS, and print TLS version, Ciphersuite name, and Issuer organization from that web
