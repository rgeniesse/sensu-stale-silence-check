package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func basicAuth() string {
	var username string = "admin"
	var passwd string = "P@ssw0rd!"
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://172.28.128.14:8080/auth", nil)
	req.SetBasicAuth(username, passwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%T\n", bodyText)
	s := string(bodyText)
	return s
}

func querySilenced() {

	// req, err := http.NewRequest("GET", "http://172.28.128.14:8080/api/core/v2/namespaces/default/silenced", nil)
}

func main() {
	fmt.Println("requesting...")
	S := basicAuth()
	fmt.Println(S)
}
