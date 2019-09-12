package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Auth struct {
	Access_token  string
	Expires_at    int
	Refresh_token string
}

func getAuthToken() string {

	myauth := Auth{}
	username := "admin"
	passwd := "P@ssw0rd!"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "http://172.28.128.14:8080/auth", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bleh := json.NewDecoder(resp.Body).Decode(&myauth)

	if bleh != nil {
		panic(bleh)
	}
	return myauth.Access_token
}

func querySilenced(token string) {
	fmt.Println(token)
	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", "http://172.28.128.14:8080/api/core/v2/namespaces/default/silenced", nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(body)))
}

func main() {
	fmt.Println("requesting...")
	token := getAuthToken()
	querySilenced(token)
	// fmt.Println("Scooby")
	// fmt.Println(S)
}
