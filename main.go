package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	Network           = "localhost"
	Port              = 9090
	MAX_PAYLOAD_IN_MB = 500_000_000
)

var Addr = fmt.Sprintf("%s:%d", Network, Port)

var SupportedBanks = map[string]bool{
	"monzo":    true,
	"halifax":  true,
	"barclays": true,
}
var AcceptedFileTypes = map[string]bool{
	// "image/png":  true,
	"text/csv": true,
	// "image/jpeg": true,
	"text": true,
}

type File struct {
	file     multipart.File
	name     string
	bankType string
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		log.Println("preflight request made");
		log.Printf("%+v\n", r);
		return
	}
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "index.html")
		return
	}

	if r.ContentLength >= MAX_PAYLOAD_IN_MB {
		log.Println("client sent a file toolarge:", r.ContentLength)
		return
	}

	userFile, headers, err := r.FormFile("user_file")
	if err != nil {
		log.Println("error extracting user_file from form")
		log.Println(err)
		http.ServeFile(w, r, "index.html")
		return
	}

	fileType := headers.Header.Get("content-type")
	if fileType == "" {
		log.Println("client did not include file-type in MIME headers")
		http.ServeFile(w, r, "index.html")
		return
	}

	log.Println("client sent fileType of: ", fileType)

	fileName := headers.Filename
	bankType := r.FormValue("bank_type")

	fileName = strings.ReplaceAll(fileName, " ", "")
	if len(fileName) == 0 || !SupportedBanks[bankType] {
		log.Println("Bad Request, either bank type isn't supported or no file was provided")
		log.Printf("filename: %s | bankType: %s\n", fileName, bankType)
		http.ServeFile(w, r, "index.html")
		return
	}

	f := &File{
		file:     userFile,
		name:     fileName,
		bankType: bankType,
	}

	res, err := f.handleMonzoBankStatement()
	if err != nil {
		log.Println(err)
		http.ServeFile(w, r, "index.html")
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("failed to encode response to send -> \n %s\n", err)
		return
	}
}

func main() {
	// http.HandleFunc("/", serverHome)
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	http.HandleFunc("/upload", handleUpload)
	log.Printf("server running at http://%s\n", Addr)
	if err := http.ListenAndServe(Addr, nil); err != nil {
		log.Fatalf(" could not start server%s", err)
	} else {
	}
}
