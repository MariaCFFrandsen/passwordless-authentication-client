package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	crypto "crypto/x509"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Load the index.html template.
var tmpl = template.Must(template.New("tmpl").ParseFiles("index.html"))

func main() {
	fmt.Println("Calling API...")
	// Serve / with the index.html file.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Serve /callme with a text response.
	http.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("You called me!")
		get()
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("pong")
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
	Username string
	PK []byte
}

func post() rsa.PublicKey {
	fmt.Println("2. Performing Http Post...")
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

	// Convert response body to Todo struct
	//var todoStruct User
	//json.Unmarshal(bodyBytes, &todoStruct)
	//	fmt.Printf("%+v\n", todoStruct)
	return privateKey.PublicKey
}

func authenticate(pk rsa.PublicKey) {
	fmt.Println("2. Performing Authenticate...")
	key := crypto.MarshalPKCS1PublicKey(&pk)
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

	// Convert response body to Todo struct
	//var todoStruct User
	//json.Unmarshal(bodyBytes, &todoStruct)
	//	fmt.Printf("%+v\n", todoStruct)
}
