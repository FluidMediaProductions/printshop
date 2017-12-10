package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"os"
)

const baseName = "hoodie"
const sourceFile = "source.png"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	baseName := r.FormValue("base")
	if baseName == "" {
		fmt.Println("No base image name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	base, err := loadImageConfig(baseName)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = processImage(base, file, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/").HandlerFunc(handleImage)
	log.Fatal(http.ListenAndServe(":8080", router))
}
