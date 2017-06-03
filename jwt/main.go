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
	"strings"

	"github.com/bmizerany/pat"

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
	m := pat.New()
	m.Post("/login", http.HandlerFunc(loginHandler))
	m.Get("/me", http.HandlerFunc(profileHandler))
	http.ListenAndServe(":3000", m)
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
type AuthToken struct {
	Username string `json:"username"`
}

func generateAuthToken(privateKey []byte, user *User) (string, error) {
	sigingKey, err := LoadPrivateKey(privateKey)
	if err != nil {
		log.Println("Error loading private key:", err)
		return "", fmt.Errorf("error loading private key")
	}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: sigingKey}, nil)

	if err != nil {
		log.Println("Error creating signer:", err)
		return "", fmt.Errorf("error creating the signer")
	}

	authToken := AuthToken{Username: user.Username}
	authTokenBytes, err := json.Marshal(authToken)
	if err != nil {
		log.Println("Error converting authToken to json:", err)
		return "", fmt.Errorf("error converting auth token to json")
	}

	objToken, err := signer.Sign(authTokenBytes)
	if err != nil {
		log.Println("Error signing the payload:", err)
		return "", fmt.Errorf("error signing the payload")
	}

	token, err := objToken.CompactSerialize()
	if err != nil {
		log.Println("Error creating compact token:", err)
		return "", fmt.Errorf("error create auth token")
	}
	return token, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if user.Username != "peterg" || user.Password != "secretone" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "{%q: %q}", "error", "invalid credentials")
		return
	}

	token, err := generateAuthToken(privateKey, user)
	log.Println("Token CompactSerialized:", token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{%q: %q}", "token", token)
}

func Authenticate() {

}

func verifyTokenFromRequest(req *http.Request, publicKey []byte) (*AuthToken, error) {
	verificationKey, err := LoadPublicKey(publicKey)
	if err != nil {
		log.Println("error loading public key", err)
		return nil, err
	}
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" || len(strings.Split(authHeader, " ")) != 2 {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := strings.Split(authHeader, " ")[1]
	parsedObj, err := jose.ParseSigned(token)
	if err != nil {
		log.Println("error parsin signed:", err)
		return nil, fmt.Errorf("error parsing signed token")
	}

	plainText, err := parsedObj.Verify(verificationKey)
	if err != nil {
		log.Println("error verifying the object:", err)
		return nil, fmt.Errorf("error verifying the object")
	}

	authToken := new(AuthToken)
	err = json.NewDecoder(bytes.NewReader(plainText)).Decode(authToken)
	return authToken, err
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authToken, err := verifyTokenFromRequest(r, publicKey)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "{%q: %q}", "error", "forbidden access")
		return
	}

	user := User{Username: authToken.Username}
	json.NewEncoder(w).Encode(user)
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
