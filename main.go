package main

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"net/http"
	"os"
	"text/template"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type Page struct {
	Title string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generator/", viewCodeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local development
	}

	println("Server starting on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "QR Code Generator"}

	t, err := template.ParseFiles("templates/generator.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, p); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func viewCodeHandler(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString")

	qrCode, err := qr.Encode(dataString, qr.L, qr.Auto)
	if err != nil {
		http.Error(w, "Error generating QR code", http.StatusInternalServerError)
		return
	}

	qrCode, err = barcode.Scale(qrCode, 512, 512)
	if err != nil {
		http.Error(w, "Error scaling QR code", http.StatusInternalServerError)
		return
	}

	// Create a buffer to store the PNG image
	var buf bytes.Buffer
	if err := png.Encode(&buf, qrCode); err != nil {
		http.Error(w, "Error encoding QR code", http.StatusInternalServerError)
		return
	}

	// Encode the image as base64
	qrBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Create template data
	data := struct {
		QRCode string
	}{
		QRCode: qrBase64,
	}

	// Parse and execute the template
	t, err := template.ParseFiles("templates/qr.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
