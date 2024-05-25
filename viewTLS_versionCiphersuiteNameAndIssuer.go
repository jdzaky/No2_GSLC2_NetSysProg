package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	var website string
	for {
		fmt.Print("Masukan URL (exit untuk keluar):")
		fmt.Scanln(&website)
		if strings.EqualFold(website, "exit") {
			fmt.Println("anda telah keluar dari program")
			os.Exit(0)
			break
		} else {
			cek(website)
		}
	}
}

func cek(website string) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 30 * time.Second},
		"tcp",
		website+":443",
		&tls.Config{
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion:       tls.VersionTLS12,
		},
	)
	if err != nil {
		fmt.Println(err)
	}

	state := conn.ConnectionState()

	fmt.Printf("TLS Version: %x\n", state.Version)

	cipherSuite := tls.CipherSuiteName(state.CipherSuite)
	fmt.Printf("Ciphersuite name: %s\n", cipherSuite)

	for _, cert := range state.PeerCertificates {
		issuerOrg := cert.Issuer.Organization
		fmt.Printf("Issuer organization: %s\n", issuerOrg)
		break
	}
}
