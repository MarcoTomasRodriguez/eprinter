package email

import (
	"encoding/json"
	"io/ioutil"

	"github.com/MarcoTomasRodriguez/eprinter/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ReadCredentials retrieves the oauth2 config from the credentials file.
func ReadCredentials(filename string) (*oauth2.Config, error) {
	// Read credentials file.
	credentials, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Load oauth2 config from credentials file.
	return google.ConfigFromJSON(credentials, config.AccessScopes...)
}

// GetToken gets a token instance from the token file
func ReadToken(filename string) (*oauth2.Token, error) {
	token := new(oauth2.Token)

	// Read the token file.
	tokenJson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse json token as oauth2.Token.
	err = json.Unmarshal(tokenJson, token)

	return token, err
}

// SaveToken saves a token into a file
func SaveToken(filename string, token *oauth2.Token) error {
	// Encode oauth2.Token as json.
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return err
	}

	// Write the token file.
	return ioutil.WriteFile(filename, tokenJson, 0600)
}
