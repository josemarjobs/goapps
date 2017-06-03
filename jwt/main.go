package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jose "gopkg.in/square/go-jose.v2"
)

var (
	privateKeyPath = "key.pem"
	publicKeyPath  = "cert.pem"
	privateKey     []byte
	publicKey      []byte
	algorithm      = "RS256"
)

func init() {
	privateKey, _ = ioutil.ReadFile(privateKeyPath)
	publicKey, _ = ioutil.ReadFile(publicKeyPath)
}

func main() {
	u := User{
		Username: "peterg",
		Name:     "Peter Griffin",
		Location: "Quahog",
	}
	sigingKey, err := LoadPrivateKey(privateKey)
	if err != nil {
		log.Println("Error loading private key:", err)
		return
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: sigingKey}, nil)
	if err != nil {
		log.Println("Error creating signer:", err)
		return
	}
	userBytes, err := json.Marshal(u)
	if err != nil {
		log.Println("error converting user to json:", err)
		return
	}
	obj, err := signer.Sign(userBytes)
	if err != nil {
		log.Println("Error signing the payload:", err)
		return
	}
	token, err := obj.CompactSerialize()
	if err != nil {
		log.Println("Error creating compact token:", err)
		return
	}
	log.Println("Token CompactSerialized:", token)

	verificationKey, err := LoadPublicKey(publicKey)
	if err != nil {
		log.Println("error loading public key", err)
		return
	}

	token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJ1c2VybmFtZSI6InBldGVyZyIsInBhc3N3b3JkIjoiIiwibmFtZSI6IlBldGVyIEdyaWZmaW4iLCJsb2NhdGlvbiI6IlF1YWhvZyJ9.1kUD0upoN9KNGWQayRFNPBex21u7expg2z57zTfAdOYCnkM9h5dwAk0QJ3WyiWEsG_m0Y2jkHMVEhUhSYwXniVtY3dYzs1TMBtsKmo1-ANdlgwnY4H-xYdBBsMqWsKcNJaf-75Hz7Vp4-ByoF85HPtXwrq4-veRY1ez5wN_MTL8NjDw4lB1R5rH-2FR6sd0YizDpBr8o_jqOhWqPgLjojElkUVNgIq1-Lpgd-QE96CPTNy5wJjIKjUiqWyZljgspEnPJGA8jR6H5bmzChikVJMDhKSYR1llQTFntS4EbjY5ZVWfaxvDPxlkj7t9MZmgNQuG1rjvlmpE0_2x4zcLUDw"
	parsedObj, err := jose.ParseSigned(token)
	if err != nil {
		log.Println("error parsed signed:", err)
		return
	}
	plainText, err := parsedObj.Verify(verificationKey)
	if err != nil {
		log.Println("error verifying:", err)
		return
	}
	log.Println("Plain text:", string(plainText))
	newUser := new(User)
	json.NewDecoder(bytes.NewReader(plainText)).Decode(newUser)
	log.Printf("New User %+v\n", newUser)

	// m := pat.New()
	// m.Post("/login", http.HandlerFunc(loginHandler))
	// m.Get("/me", http.HandlerFunc(profileHandler))
	// http.ListenAndServe(":3000", m)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{token: %q}", "token")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

}

// LoadPublicKey loads a public key from PEM/DER-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

// LoadPrivateKey loads a private key from PEM/DER-encoded data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s' and '%s'", err0, err1, err2)
}
