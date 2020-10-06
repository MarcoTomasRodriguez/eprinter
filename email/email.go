package email

import (
	"github.com/MarcoTomasRodriguez/auto-printer/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
)

// GetConfiguration retrieves the configuration from the credentials file.
func GetConfiguration() (*oauth2.Config, error) {
	credentials, err := ioutil.ReadFile(config.CredentialsFilename)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(credentials, config.AccessScopes...)
}

// GetEmailById gets the email by its id
func GetEmailById(mailId string, service *gmail.Service) (*gmail.Message, error) {
	return service.Users.Messages.Get("me", mailId).Do()
}
