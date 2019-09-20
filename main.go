// TODO:
// Comment code

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
	err2 := json.NewDecoder(resp.Body).Decode(&myauth)

	if err2 != nil {
		log.Fatal(err2)
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
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err2 := json.Unmarshal([]byte(body), silenced2)

	if err2 != nil {
		log.Fatal(err)
	}
}

func checkIfSilencedOld(silenced3 []Silenced) {

	if len(silenced3) > 0 {
		for i := 0; i < len(silenced3); i++ {
			n := time.Unix(int64(silenced3[i].Begin), 0)
			duration := time.Since(n)
			not_too_old := "An entry for "

			fmt.Println(n)

			if int(duration.Seconds()) > threshold && silenced3[i].Expire == int(-1) && !silenced3[i].Expire_on_resolve {
				fmt.Println("A silenced entry " + silenced3[i].Metadata.Name + " has been active since " + n.String())
			}

			// Handle cases where the silenced entry should not be flagged
			// To Fix: not getting into these else if statements.

			if int(duration.Seconds()) < threshold && not_too_old == "An entry for " {
				fmt.Println(not_too_old + silenced3[i].Metadata.Name + " was not flagged because the threshold of not met")
			} else if int(duration.Seconds()) < threshold {

			}

			if silenced3[i].Expire != int(-1) && not_too_old == "An entry for " {
				fmt.Println(not_too_old + silenced3[i].Metadata.Name + " was not flagged as the silence is set to expire after some time")
			} else if silenced3[i].Expire != int(-1) {
				fmt.Println("you will expire after sometime as well")

			}

			if silenced3[i].Expire_on_resolve && not_too_old == "An entry for " {
				fmt.Println(not_too_old + silenced3[i].Metadata.Name + " was not flagged as the silence is to to expire once the  ")
			} else if silenced3[i].Expire_on_resolve {
				fmt.Println("you will expire on resolve as well")

			}

		}
		os.Exit(1)
	} else {
		fmt.Println("Good news nobody, the silenced endpoint is empty!")
		os.Exit(0)
	}
}

func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}

	token := getAuthToken()
	silenced := []Silenced{}
	querySilenced(token, &silenced)
	checkIfSilencedOld(silenced)

}
