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

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Silenced struct {
	Metadata          Metadata
	Expire            int    `json:"expire"`
	Expire_on_resolve bool   `json:"expire_on_resolve"`
	Creator           string `json:"creator"`
	Check             string `json:"check"`
	Subscription      string `json:"subscription"`
	Begin             int    `json:"begin"`
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

	// fmt.Printf("%v\n", myauth)
	// fmt.Printf("%T\n", myauth)
	return myauth.Access_token
}

func querySilenced(token string) {

	silenced := []Silenced{}
	bearer := "Bearer " + token
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "http://172.28.128.14:8080/api/core/v2/namespaces/default/silenced", nil)
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bleh := json.Unmarshal([]byte(body), &silenced)

	if bleh != nil {
		panic(bleh)
	}
	// fmt.Printf("%v\n", silenced)
	fmt.Printf("%v\n", silenced[0].Begin)
	fmt.Printf("%v\n", silenced[0].Check)
	fmt.Printf("%v\n", silenced[0].Creator)
	fmt.Printf("%v\n", silenced[0].Expire)
	fmt.Printf("%v\n", silenced[0].Expire_on_resolve)
	fmt.Printf("%v\n", silenced[0].Subscription)
	fmt.Printf("%v\n", silenced[0].Metadata.Name)
	fmt.Printf("%v\n", silenced[0].Metadata.Namespace)
}

func main() {
	fmt.Println("requesting...")
	token := getAuthToken()
	querySilenced(token)
}
