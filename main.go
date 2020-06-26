package main

import (
	"github.com/MarcoTomasRodriguez/auto-printer/auth"
	"github.com/MarcoTomasRodriguez/auto-printer/email"
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"os"
)

func runService() {
	// Loads the context
	ctx := context.Background()

	// Gets the configuration from the credentials file
	configuration, err := auth.GetConfiguration()
	if err != nil {
		log.Fatalf("Unable to parse client secret fs to conf: %v", err)
	}

	// Gets the token
	token := auth.GetToken()

	// Gets the service instance
	service, err := gmail.NewService(ctx, option.WithTokenSource(configuration.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// Retrieves printed identifier label id
	printedLabelId, err := email.GetPrintedLabelId(service)
	if err != nil {
		log.Fatalf("Unable to get printed label id: %v", err)
	}

	// Returns all the unprinted emails
	res, err := email.GetUnlabelledEmails(service)
	if err != nil {
		log.Fatalf("Unable to retrieve emails: %v", err)
	}

	// For each email
	for _, mail := range res.Messages {
		mailId := mail.Id

		if mail == nil {
			log.Fatalf("Mail should never be nil")
		}

		// Gets more information about the email
		mail, err = email.GetEmailById(mailId, service)
		if err != nil {
			log.Fatalf("Unable to get mail %v", err)
		}

		// Gets its attachments
		attachments, err := email.GetAttachments(mail, service)
		if err != nil {
			log.Fatalf("Unable to get attachment %v", err)
		}

		// For each attachment
		for _, attachment := range attachments {
			attachmentId := attachment.Body.AttachmentId

			if attachment == nil {
				log.Fatalf("Attachment should never be nil")
			}

			if attachment.PartId == "0" {
				continue
			}

			// Gets the attachment body
			attachmentBody, err := email.GetAttachmentBody(attachmentId, mailId, service)
			if err != nil {
				log.Fatalf("Unable to get attachment body %v", err)
			}

			// Prints the attachment
			if err = email.PrintAttachment(attachment, attachmentBody); err != nil {
				log.Fatalf("Unable to print attachment %v", err)
			}

			log.Printf("The attachment was successfully printed. Mail ID: %s - Attachment ID: %s", mailId, attachmentId)
		}

		// Sets the email as printed
		if _, err = email.SetPrintedEmail(printedLabelId, mailId, service); err != nil {
			log.Fatalf("Unable to add printed label to email %v", err)
		}

		log.Printf("The email was successfully labeled as printed. Mail ID: %s", mailId)
	}
}

func runAuth() {
	// Gets the configuration from the credentials file
	configuration, err := auth.GetConfiguration()
	if err != nil {
		log.Fatalf("Unable to parse client secret to configuration: %v", err)
	}

	// Updates the token
	auth.UpdateToken(configuration)
}

func printUsage() {
	log.Printf("Usage: %s {service|auth}", os.Args[0])
	os.Exit(1)
}

func main() {
	args := os.Args

	if len(args) < 2 {
		printUsage()
	}

	switch args[1] {
	case "service":
		runService()
	case "auth":
		runAuth()
	default:
		printUsage()
	}

}