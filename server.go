package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	tlsConfig := &tls.Config{
		CurvePreferences:         []tls.CurveID{tls.CurveP256},
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
	}

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	listener, err := tls.Listen("tcp", ":8080", tlsConfig)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/json", handleJSON)
	http.HandleFunc("/upload", handleUpload)
	fmt.Println("Server started on http://localhost:8080")
	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Received JSON data: %v", data)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// post
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Simpan file ke direktori uploads/
		uploadedFilePath := fmt.Sprintf("./uploads/%s", handler.Filename)
		data, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = ioutil.WriteFile(uploadedFilePath, data, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File Uploaded: %s", handler.Filename)

		// fmt.Fprintf(w, "Received file data: %s", data)
	case http.MethodGet:
		// GET
		fmt.Fprint(w, `
			<html>
			<body>
				<form action="/upload" method="post" enctype="multipart/form-data">
					<input type="file" name="file">
					<input type="submit" value="Upload">
				</form>
			</body>
			</html>
		`)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
