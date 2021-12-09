package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/MarcoTomasRodriguez/eprinter/config"
	"github.com/MarcoTomasRodriguez/eprinter/email"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// setupCmd represents the setup command.
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Loads the credentials and generates an oauth2 access token.",
	Long:  "Loads the credentials and generates an oauth2 access token.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the credentials file.
		credentials, err := email.ReadCredentials(args[0])
		if err != nil {
			log.Fatalf("Could not load configuration: %v", err)
		}

		// Redirect user to the auth consent page.
		authConsentUrl := credentials.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		log.Printf("Go to the following link in your browser: %s\nEnter the authorization code: ", authConsentUrl)

		// Input authorization code.
		var authorizationCode string
		if _, err := fmt.Scan(&authorizationCode); err != nil {
			log.Fatalf("Unable to read authorization code: %v", err)
		}

		// Convert the authorization code into a token.
		token, err := credentials.Exchange(context.Background(), authorizationCode)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web: %v", err)
		}

		// Copy credentials file.
		credentialsFile, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatalf("Could not read configuration file: %v", err)
		}
		if err := ioutil.WriteFile(config.Config.CredentialsFilename(), credentialsFile, 0600); err != nil {
			log.Fatalf("Could not write configuration file: %v", err)
		}

		// Save the oauth2 token.
		email.SaveToken(config.Config.TokenFilename(), token)

		log.Println("Setup completed successfully.")
	},
}

// init registers the setup command.
func init() {
	rootCmd.AddCommand(setupCmd)
}
