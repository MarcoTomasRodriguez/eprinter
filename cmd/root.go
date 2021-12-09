package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/MarcoTomasRodriguez/eprinter/config"
	"github.com/MarcoTomasRodriguez/eprinter/email"
	"github.com/MarcoTomasRodriguez/eprinter/printer"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "eprinter",
	Short: "Automatically print email attachments.",
	Long:  "Automatically print email attachments.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Read the credentials.
		credentials, err := email.ReadCredentials(config.Config.CredentialsFilename())
		if err != nil {
			log.Fatalf("Could not load credentials: %v", err)
		}

		// Read the token.
		token, err := email.ReadToken(config.Config.TokenFilename())
		if err != nil {
			log.Fatalf("Could not load token: %v", err)
		}

		// Create gmail service wrapper.
		emailService, err := email.NewEmailService(ctx, credentials, token)
		if err != nil {
			log.Fatalf("Could not create gmail client: %v", err)
		}

		// Find or create printed label.
		printedLabel, err := emailService.GetLabelByName(config.Config.PrintedLabel)
		if err != nil && err.Error() == "no labels match the desired criteria" {
			printedLabel, err = emailService.CreateLabel(config.Config.PrintedLabel)
			if err != nil {
				log.Fatalf("Could not create printed label: %v", err)
			}
		}

		// Emails should have attachments and no user labels.
		query := "has:attachment -(label:" + printedLabel.Name + ")"

		// Add emails to query.
		for _, email := range config.Config.AllowedEmails {
			query += " from:" + email
		}

		// Add subjects to query.
		for _, subject := range config.Config.AllowedEmailSubjects {
			query += " subject:" + subject
		}

		// Find mails that match the query.
		mails, err := emailService.FindEmails(query)
		if err != nil {
			log.Fatalf("Could not find emails: %v", err)
		}

		// For each email
		for _, mail := range mails {
			// Get complete information about the email.
			mail, err = emailService.GetEmailById(mail.Id)
			if err != nil {
				log.Fatalf("Could not find email by id: %v", err)
			}

			// Iterate over email attachments.
			for _, attachment := range mail.Payload.Parts {
				// Ignore attachment 0.
				// It is an internal attachment that does not contain meaningful data for this program.
				if attachment.PartId == "0" {
					continue
				}

				// Save the attachment.
				if err := emailService.SaveAttachment(mail, attachment); err != nil {
					log.Fatalf("Could not save attachment: %v", err)
				}

				// Print the attachment.
				if err := printer.PrintFile(attachment.Filename); err != nil {
					log.Fatalf("Could not print file: %v", err)
				}

				log.Printf("Attachment successfully printed - Mail ID: %s - Attachment ID: %s", mail.Id, attachment.Body.AttachmentId)
			}

			// Label email as printed.
			if _, err = emailService.SetEmailLabel(mail.Id, printedLabel.Id); err != nil {
				log.Fatalf("Could not label email as printed: %v", err)
			}

			log.Printf("Email labeled as printed - Mail ID: %s", mail.Id)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Define config global flag.
	homeDir, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().StringVar(&config.Filename, "config", filepath.Join(homeDir, ".eprinter/config.toml"), "Set config file.")

	// Initialize config.
	cobra.OnInitialize(config.LoadConfig)
}
