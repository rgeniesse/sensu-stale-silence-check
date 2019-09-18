// Todo:
// Add data structure to hold all endpoints found to be too old.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
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

var (
	username, password, host, port string
	threshold, timeout             int
)

func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-stale-silence-check",
		Short: "A Sensu Go check plugin to send out reminders about stale silenced entries",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&username,
		"username",
		"u",
		os.Getenv("SENSU_API_USER"),
		"A Sensu Go user with API access.")
	cmd.MarkFlagRequired("username")

	cmd.Flags().StringVarP(&password,
		"password",
		"p",
		os.Getenv("SENSU_API_PASSWORD"),
		"A Sensu Go user's password.")
	cmd.MarkFlagRequired("password")

	cmd.Flags().StringVarP(&host,
		"host",
		"H",
		os.Getenv("SENSU_API_HOST"),
		"The Sensu API host.")
	cmd.MarkFlagRequired("host")

	cmd.Flags().StringVarP(&port,
		"port",
		"P",
		"8080",
		"The port the Sensu API is listening on.")

	cmd.Flags().IntVarP(&threshold,
		"threshold",
		"t",
		604800,
		"Threshold in seconds to consider a silenced entry stale")

	cmd.Flags().IntVarP(&timeout,
		"timeout",
		"T",
		10,
		"Time in seconds to consider the API unresponsive")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	return nil
}

func getAuthToken() string {

	myauth := Auth{}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/auth", nil)
	req.SetBasicAuth(username, password)
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

func querySilenced(token string, silenced2 *[]Silenced) {

	bearer := "Bearer " + token
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/api/core/v2/namespaces/default/silenced", nil)
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

func checkIfSilencedOld(silenced3 []Silenced) {

	fmt.Println(silenced3[0].Begin)
	fmt.Printf("%T\n", silenced3)
	fmt.Println(silenced3)

	if len(silenced3) > 0 {
		for i := 0; i < len(silenced3); i++ {
			n := time.Unix(int64(silenced3[i].Begin), 0)
			duration := time.Since(n)

			fmt.Println(silenced3[i])
			if int(duration.Seconds()) > threshold {
				// Make a data structure that holds all the entries. Use that data structure
				// to determine the exit status, otherwise program exits on first hit.
				fmt.Println("This entry is old and was added to check result!")
			} else {
				fmt.Println("Entry was not added to check result")
			}

		}
		os.Exit(1)
	} else {
		fmt.Println("Silenced endpoint is empty!")
		os.Exit(0)
	}
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Requesting Auth Token")
	token := getAuthToken()
	silenced := []Silenced{}
	fmt.Printf("%T\n", silenced)
	querySilenced(token, &silenced)
	checkIfSilencedOld(silenced)

}
