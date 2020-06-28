package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
)

const credentialsFilename = "credentials.json"

var scopes = []string{gmail.GmailLabelsScope, gmail.GmailReadonlyScope, gmail.GmailModifyScope}

// Retrieves the configuration from the credentials file.
func GetConfiguration() (*oauth2.Config, error) {
	credentials, err := ioutil.ReadFile(credentialsFilename)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(credentials, scopes...)
}


