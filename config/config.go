package config

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/pelletier/go-toml"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ProgramFolder is the folder where the program (and the user) stores his information
var ProgramFolder = filepath.Join(os.Getenv("HOME"), ".auto-printer")

// ConfigFilename is the editable configuration file
var ConfigFilename = filepath.Join(ProgramFolder, "config.toml")

// TokenFilename is the token file generated when auth occurs
var TokenFilename = filepath.Join(ProgramFolder, "token.json")

// CredentialsFilename is the credentials file downloaded by the user required for the program to work
var CredentialsFilename = filepath.Join(ProgramFolder, "credentials.json")

// AccessScopes are the different Gmail API permissions granted by the user
var AccessScopes = []string{gmail.GmailLabelsScope, gmail.GmailReadonlyScope, gmail.GmailModifyScope}


// Config is the editable configuration struct
type Config struct {
	AllowedEmails        []string `toml:"allowed_emails"`
	AllowedEmailSubjects []string `toml:"allowed_email_subjects"`
	PrintedLabelName     string   `toml:"printed_label_name"`
}

// LoadConfig loads the configuration
func LoadConfig() *Config {
	config := &Config{}

	file, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	if err = toml.Unmarshal(file, config); err != nil {
		log.Fatalf("Unable to parse configuration file: %v", err)
	}

	return config
}

// ValidateAllowedEmails checks whether the allowed emails field in the configuration is valid
func (config *Config) ValidateAllowedEmails() error {
	if len(config.AllowedEmails) == 0 {
		return errors.New("allowed emails should never be empty")
	}

	for _, allowedEmail := range config.AllowedEmails {
		if err := checkmail.ValidateFormat(allowedEmail); err != nil {
			return err
		}
	}

	return nil
}

// ValidateAllowedEmailSubjects checks whether the allowed email subjects field in the configuration is valid
func (config *Config) ValidateAllowedEmailSubjects() error {
	if len(config.AllowedEmailSubjects) == 0 {
		return errors.New("allowed email subjects should never be empty")
	}

	for _, allowedEmailTitle := range config.AllowedEmailSubjects {
		if len(strings.TrimSpace(allowedEmailTitle)) == 0 {
			return errors.New("allowed email subject should never be empty")
		}
		if len(strings.Split(allowedEmailTitle, " ")) != 1 {
			return errors.New("allowed email subject should have only one word")
		}
	}

	return nil
}

// ValidatePrintedLabelName checks whether the printed label name in the configuration is valid
func (config *Config) ValidatePrintedLabelName() error {
	if len(strings.TrimSpace(config.PrintedLabelName)) == 0 {
		return errors.New("printed label name should never be empty")
	}

	return nil
}


