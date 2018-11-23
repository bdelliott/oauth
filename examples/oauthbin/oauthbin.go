package main

import (
	"fmt"
	"log"

	"github.com/mrjones/oauth"
)


func main() {
	consumerKey := "key"
	consumerSecret := "secret"

	c := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "http://oauthbin.com/v1/request-token",
			AuthorizeTokenUrl: "http://oauthbin.com/v1/authorize",
			AccessTokenUrl:    "http://oauthbin.com/v1/access-token",
			HttpMethod:        "GET",
		})

	c.Debug(true)

	requestToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("(1) Go to: " + url)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(accessToken)
	/*
	client, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Get(*jiraUrl + "/rest/api/2/issue/BULK-1")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	fmt.Println("Data: " + string(bits))
*/
}
