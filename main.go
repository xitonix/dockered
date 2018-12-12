package main // import "go.xitonix.io/dockered"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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

	encrypted, err := encryptData([]byte(data))
	if err != nil {
		http.Error(w, "Failed to encrypt the data", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(encrypted))
	if err != nil {
		http.Error(w, "Failed to write the response", http.StatusInternalServerError)
		return
	}
}

func encryptData(input []byte) (string, error) {

	// ** ATTENTION **
	// NEVER EVER EVER commit the encryption key into the source control
	// This is for demonstration purposes ONLY
	key := []byte("16 bytes AES key")

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(input))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], input)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}
