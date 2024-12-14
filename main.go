package main

import (
	"image/png"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"fmt"
	"os"
	"text/template"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/spf13/viper"
)

type QrText struct {
	Text string `json:"text"`
}

type Page struct {
	Title string
}

func main() {

	viper.SetConfigFile("ENV")
	viper.ReadInConfig()
	viper.AutomaticEnv()
	port := fmt.Sprint(viper.Get("PORT"))

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/generator/", qrView).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	log.Println(http.ListenAndServe(":"+port, loggedRouter))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "QR Code Generator"}

	t, err := template.ParseFiles("generator.html")
	if err != nil {
		log.Println("Problem parsing html file")
	}

	t.Execute(w, p)
}

func qrView(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString")

	qrCode, err := qr.Encode(dataString, qr.L, qr.Auto)
	if err != nil {
		fmt.Println(err)
	} else {
		qrCode, err = barcode.Scale(qrCode, 128, 128)
		if err != nil {
			fmt.Println(err)
		} else {
			png.Encode(w, qrCode)
		}
	}

}
