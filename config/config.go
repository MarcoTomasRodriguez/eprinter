package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/spf13/viper"
	"google.golang.org/api/gmail/v1"
)

// Configuration defines the behavior of the application.
type Configuration struct {
	// ProgramFolder is the folder used to store program data.
	ProgramFolder string `mapstructure:"program_folder"`

	// AllowedEmails are the sender emails that the program will accept.
	AllowedEmails []string `mapstructure:"allowed_emails"`

	// AllowedEmailSubjects are the subjects that the program will accept.
	AllowedEmailSubjects []string `mapstructure:"allowed_email_subjects"`

	// PrintedLabel is the label used to mark emails as printed.
	PrintedLabel string `mapstructure:"printed_label"`
}

// Validate validates the config.
func (config *Configuration) Validate() error {
	// Validate AllowedEmails.
	if len(config.AllowedEmails) == 0 {
		return errors.New("allowed emails should never be empty")
	}
	for _, allowedEmail := range config.AllowedEmails {
		if err := checkmail.ValidateFormat(allowedEmail); err != nil {
			return err
		}
	}

	// Validate AllowedEmailSubjects.
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

	// Validate PrintedLabelName.
	if len(strings.TrimSpace(config.PrintedLabel)) == 0 {
		return errors.New("printed label name should never be empty")
	}

	return nil
}

// TokenFilename returns the path to the token file.
func (config *Configuration) TokenFilename() string {
	return filepath.Join(config.ProgramFolder, "token.json")
}

// CredentialsFilename returns the path to the credentials file.
func (config *Configuration) CredentialsFilename() string {
	return filepath.Join(config.ProgramFolder, "credentials.json")
}

// AccessScopes are the different Gmail API permissions granted by the user
var AccessScopes = []string{gmail.GmailLabelsScope, gmail.GmailReadonlyScope, gmail.GmailModifyScope}

// Filename is the path of the configuration file.
var Filename string

// Config is the shared configuration instance.
var Config = &Configuration{}

// LoadConfig loads the configuration
func LoadConfig() {
	// Check if config file exists.
	if _, err := os.Stat(Filename); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: unable to open configuration file: %v\n", err)
		panic(err)
	}

	// Set config file.
	viper.SetConfigFile(Filename)

	// Read config.
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: unable to read configuration file: %v\n", err)
		panic(err)
	}

	// Get home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: could not get home directory %v\n", err)
		panic(err)
	}

	// Set defaults.
	viper.SetDefault("program_folder", filepath.Join(homeDir, ".eprinter"))

	// Unmarshal configuration into the shared config struct.
	if err := viper.Unmarshal(Config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: unable to decode configuration into struct: %v\n", err)
		panic(err)
	}

	// Validate configuration.
	if err := Config.Validate(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: invalid configuration: %v\n", err)
		panic(err)
	}
}
