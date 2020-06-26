package conf

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
	"strings"
)

type Config struct {
	AllowedEmails        []string `toml:"allowed_emails"`
	AllowedEmailSubjects []string `toml:"allowed_email_subjects"`
	PrintedLabelName     string   `toml:"printed_label_name"`
}

const configFilename = "/etc/auto-printer.toml"

// Loads the configuration
func LoadConfig() *Config {
	config := &Config{}

	file, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("Unable to read configuration file: %v", err)
	}

	if err = toml.Unmarshal(file, config); err != nil {
		log.Fatalf("Unable to parse configuration file: %v", err)
	}

	return config
}

// Checks whether the allowed emails field in the configuration is valid
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

// Checks whether the allowed email subjects field in the configuration is valid
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

// Checks whether the printed label name in the configuration is valid
func (config *Config) ValidatePrintedLabelName() error {
	if len(strings.TrimSpace(config.PrintedLabelName)) == 0 {
		return errors.New("printed label name should never be empty")
	}

	return nil
}


