package main

import (
	"flag"
	"fmt"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Usage() {
	fmt.Println("Usage:")
	fmt.Print("go run examples/fatsecret/fatsecret.go")
	fmt.Print("  --consumerkey <consumerkey>")
	fmt.Println("  --consumersecret <consumersecret>")
	fmt.Println("")
	fmt.Println("In order to get your consumer key and consumer secret, you must register an 'app' at:")
	fmt.Println("https://platform.fatsecret.com/api/")
}

func main() {
	var consumerKey *string = flag.String(
		"consumerkey",
		"",
		"Consumer Key from FatSecret.")

	var consumerSecret *string = flag.String(
		"consumersecret",
		"",
		"Consumer Secret from FatSecret.")

	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 {
		fmt.Println("You must set the --consumerkey and --consumersecret flags.")
		fmt.Println("---")
		Usage()
		os.Exit(1)
	}

	c := oauth.NewConsumer(
		*consumerKey,
		*consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "http://www.fatsecret.com/oauth/request_token",
			AuthorizeTokenUrl: "http://www.fatsecret.com/oauth/authorize",
			AccessTokenUrl:    "http://www.fatsecret.com/oauth/access_token",
			ParamsInURI: true, // The FatSecret API fails unless the auth params are in the Request URI
		})

	c.Debug(true)

	requestToken, u, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("(1) Go to: " + u)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://platform.fatsecret.com/rest/server.api"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("method", "foods.search")
	q.Add("format", "json")
	q.Add("search_expression", "pizza")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Search results:")
	fmt.Println(string(results))
}
