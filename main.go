package main // import "go.xitonix.io/dockered"

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/encrypt/{data}", encrypt)
	const port = ":8080"
	fmt.Printf("The server is running on port %s\n", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func encrypt(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, ok := vars["data"]
	if !ok || len(data) == 0 {
		http.Error(w, "Data is not provided", http.StatusBadRequest)
		return
	}

	hs := sha512.New()
	hs.Write([]byte(data))
	_, err := w.Write([]byte(hex.EncodeToString(hs.Sum(nil))))
	if err != nil {
		http.Error(w, "Failed to encrypt the data", http.StatusInternalServerError)
	}
}
