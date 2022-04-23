package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	crypto "crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/MariaCFFrandsen/passwordless-authentication/authenticator/cryptography"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Load the index.html template.
var tmpl = template.Must(template.New("tmpl").ParseFiles("frontend/index.html"))

func main() {
	fmt.Println("Ready to authenticate")
	// Serve / with the index.html file.
	http.HandleFunc("/frontend/", func(w http.ResponseWriter, r *http.Request) { //rm?
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Serve /callme with a text response.
	http.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {
		pk := post()
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		if pk == nil {
			json.NewEncoder(w).Encode("already have a user")
		} else {
			json.NewEncoder(w).Encode("authenticated")
		}

	})

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		certificate := cryptography.RetrieveCertificate("hello")
		authenticated := authenticate(certificate.PublicKey)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if authenticated {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	// Start the server at http://localhost:9000
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func get() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:4000/ping", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	//req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	//defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	//var responseObject Response
	//json.Unmarshal(bodyBytes, &responseObject)
	//fmt.Printf("API Response as struct %+v\n", responseObject)
	s := string(bodyBytes)
	fmt.Printf("msg: %s", s)
}

type User struct {
	Username string `json:"username"`
	PK       []byte `json:"publickey"`
}
func post() *rsa.PublicKey {
	if FileExists("hello-key.txt") || FileExists("hello-certificate.bin") {
		return nil
	}

	fmt.Println("1. Creating User...")
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	key := crypto.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	user := User{
		Username: "maria",
		PK:       key,
	}
	jsonReq, err := json.Marshal(user)
	resp, err := http.Post("http://localhost:4000/create-user", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	if resp.StatusCode == 201  {
		cryptography.SaveCertificate(cryptography.CreateCertificate(cryptography.KeyPair{
			PrivateKey: &cryptography.PrivateKey{PrivateKey: privateKey},
			PublicKey:  &cryptography.PublicKey{PublicKey: &privateKey.PublicKey},
		}), "hello")
	}
	return &privateKey.PublicKey
}

func authenticate(key []byte) bool {
	fmt.Println("2. Authenticating...")
	user := User{
		Username: "maria",
		PK:       key,
	}
	jsonReq, err := json.Marshal(user)
	resp, err := http.Post("http://localhost:4000/authenticate-user", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	if resp.StatusCode == 202 {
		return true
	}
	return false
}
