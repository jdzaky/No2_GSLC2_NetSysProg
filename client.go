package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

func main() {
	// Get request
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

	// Post JSON request
	jsonData := map[string]string{"message": "Hello, World!"}
	resp, err = postJSONRequest("http://localhost:8080/json", jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response Body:", string(body))

	// Post multipart form request
	resp, err = postMultipartFormRequest("http://localhost:8080/upload", "file.txt", []byte("This is a sample file."))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response Body:", string(body))
}

func getRequest(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add cancellation mechanism
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
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
		Transport: &http.Transport{
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

	// Add cancellation mechanism
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
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
		Transport: &http.Transport{
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

	// Add cancellation mechanism
	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
