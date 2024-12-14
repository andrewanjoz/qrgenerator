package main

import (
	"image/png"
	"net/http"
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

	println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
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

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	png.Encode(w, qrCode)
}
