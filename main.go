// Todo:
// Add command line paremeter pass in threshold, API credentials, API host, and API port.
// Add looping logic for silenced end points with more than one entry.
// Add data structure to hold all endpoints found to be too old.
// Make exit status based on if that data structure has entries or not.
// Add check if silenced is empty as program crashes with empty array.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	req, err := http.NewRequest("GET", "http://172.28.128.15:8080/auth", nil)
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

func querySilenced(token string, silenced2 *[]Silenced) {

	bearer := "Bearer " + token
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "http://172.28.128.15:8080/api/core/v2/namespaces/default/silenced", nil)
	req.Header.Add("Authorization", bearer)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bleh := json.Unmarshal([]byte(body), silenced2)

	if bleh != nil {
		panic(bleh)
	}
}

func checkIfSilencedOld(t int64, time_treshold int) {
	n := time.Unix(t, 0)
	fmt.Println(n)
	duration := time.Since(n)
	fmt.Println(duration.Seconds())

	if int(duration.Seconds()) > time_treshold {
		// Make a data structure that holds all the entries. Use that data structure
		// to determine the exit status, otherwise program exits on first hit.
		fmt.Println("This entry is old and was added to check result!")
		os.Exit(1)
	} else {
		fmt.Println("Entry was not added to check result")
	}
}

func main() {
	fmt.Println("requesting...")
	token := getAuthToken()
	silenced := []Silenced{}
	querySilenced(token, &silenced)
	time_treshold := 109000

	// fmt.Printf("%v\n", silenced)
	// fmt.Printf("%v\n", silenced[0].Begin)
	// fmt.Printf("%v\n", silenced[0].Check)
	// fmt.Printf("%v\n", silenced[0].Creator)
	// fmt.Printf("%v\n", silenced[0].Expire)
	// fmt.Printf("%v\n", silenced[0].Expire_on_resolve)
	// fmt.Printf("%v\n", silenced[0].Subscription)
	// fmt.Printf("%v\n", silenced[0].Metadata.Name)
	// fmt.Printf("%v\n", silenced[0].Metadata.Namespace)

	checkIfSilencedOld(int64(silenced[0].Begin), time_treshold)

}
