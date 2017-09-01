package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

const URL = "https://api.github.com/graphql"

var token = os.Getenv("TOKEN")

func main() {
	query := `query {
	


    repository (name:"cuneiform", owner:"peernova-private") {
      pullRequests (last:10, baseRefName:"master"){
        nodes {
        url
        number
        headRefName

        }
      }
    }
  }`

	query = strings.Replace(query, "\n", " ", -1)
	query = strings.Replace(query, "\t", " ", -1)
	query = strings.Replace(query, "\"", "\\\"", -1)

	str := `{ "query" : "` + query + `"}`

	spew.Dump(str)

	req, err := http.NewRequest("POST", URL, strings.NewReader(str))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Body:", string(body))
}
