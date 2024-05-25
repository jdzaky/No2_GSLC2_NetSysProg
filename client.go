package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Get request")
		fmt.Println("2. Post JSON request")
		fmt.Println("3. Post multipart form request")
		fmt.Println("4. Exit")

		fmt.Print("Enter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			resp, err := getRequest("http://localhost:8080/")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close() // Close the response body to prevent resource leaks

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Response Body:", string(body))
		case "2":
			jsonData := map[string]string{"message": "Hello, World!"}
			resp, err := postJSONRequest("http://localhost:8080/json", jsonData)
			// resp, err := http.Get("http://localhost:8080/json")

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Response Body:", string(body))
		case "3":
			resp, err := postMultipartFormRequest("https://localhost:8080/upload", "file.txt", []byte("This is a sample file."))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Please go to http://localhost:8080/upload in your browser")
				return
			}
			fmt.Println("Response Body:", string(body))
		case "4":
			fmt.Println("thank you")
			os.Exit(0)
		default:
			fmt.Println("input not valid")
		}
	}
}

func getRequest(url string) (*http.Response, error) {
	// Baca sertifikat CA dari file cert.pem
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		CurvePreferences:   []tls.CurveID{tls.CurveP256},
		MinVersion:         tls.VersionTLS12,
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// cancellation mechanism
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func postJSONRequest(url string, data interface{}) (*http.Response, error) {
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		RootCAs:          caCertPool,
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
		Transport: &http.Transport{
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: true, // Disable persistent TCP connections
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// cancellation mechanism
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func postMultipartFormRequest(url, fieldName string, data []byte) (*http.Response, error) {
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		RootCAs:          caCertPool,
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
		Transport: &http.Transport{
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: true, // Disable persistent TCP connections
		},
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, "file.txt")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// cancellation mechanism
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
