package auth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"os"
)

// Token file location
const tokenFilename = "token.json"

// Gets a token instance from the token file
func GetToken() *oauth2.Token {
	token, err := getTokenFromFile()
	if err != nil {
		log.Fatalf("Unable to get token: %s", err)
	}

	return token
}

// Saves the token into the token file
func UpdateToken(config *oauth2.Config) {
	saveToken(getTokenFromWeb(config))
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	var err error
	var token *oauth2.Token
	var authCode string

	// Gets the auth url
	var authURL = config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	// Redirects
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	// Reads token from stdin
	if _, err = fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	// Converts the auth code into a valid token
	token, err = config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token
}

// Retrieves a token from a local fs.
func getTokenFromFile() (*oauth2.Token, error) {
	var file *os.File
	var token = &oauth2.Token{}
	var err error

	// Opens the token file
	file, err = os.Open(tokenFilename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Decodes the token
	token = &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)

	return token, err
}

// Saves a token into a file
func saveToken(token *oauth2.Token) {
	var file *os.File
	var err error

	log.Printf("Saving credentials in: %s\n", tokenFilename)

	// Opens the token file with creation perms
	file, err = os.OpenFile(tokenFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer file.Close()

	// Encodes the token
	err = json.NewEncoder(file).Encode(token)
	if err != nil {
		log.Fatalf("Unable to encode token: %v", err)
	}
}