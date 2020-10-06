package email

import (
	"fmt"
	"github.com/MarcoTomasRodriguez/auto-printer/config"
	"google.golang.org/api/gmail/v1"
)

// GetPrintedLabelId gets (if it doesn't exist, creates) the printed label, and returns his id
func GetPrintedLabelId(service *gmail.Service) (string, error) {
	// Gets the configuration and validates the used fields
	cfg := config.LoadConfig()
	if err := cfg.ValidatePrintedLabelName(); err != nil {
		return "", err
	}

	printedLabelName := cfg.PrintedLabelName

	// Gets the labels list
	labels, err := service.Users.Labels.List("me").Do()
	if err != nil {
		return "", err
	}

	// Finds the label id
	for _, label := range labels.Labels {
		if label.Name == printedLabelName {
			return label.Id, nil
		}
	}

	// Creates the label id if none
	data := &gmail.Label{Name: printedLabelName}
	label, err := service.Users.Labels.Create("me", data).Do()
	if err != nil {
		return "", err
	}

	return label.Id, nil
}

// GetUnlabelledEmails gets all the unlabelled emails which the program has to print
func GetUnlabelledEmails(service *gmail.Service) (*gmail.ListMessagesResponse, error) {
	var err error
	var allowedEmails string
	var allowedEmailSubjects string

	// Loads the configuration and validates the used fields
	cfg := config.LoadConfig()
	if err = cfg.ValidateAllowedEmails(); err != nil {
		return nil, err
	}
	if err = cfg.ValidateAllowedEmailSubjects(); err != nil {
		return nil, err
	}

	// Parses the allowed emails into a valid gmail query
	for _, allowedEmail := range cfg.AllowedEmails {
		allowedEmails = fmt.Sprintf("%s from:%s", allowedEmails, allowedEmail)
	}

	// Parses the allowed subjects into a valid gmail query
	for _, allowedEmailSubject := range cfg.AllowedEmailSubjects {
		allowedEmailSubjects = fmt.Sprintf("%s subject:%s", allowedEmailSubjects, allowedEmailSubject)
	}

	// Parses the query and executes it
	query := fmt.Sprintf("{%s} {%s} has:attachment has:nouserlabels", allowedEmails, allowedEmailSubjects)
	return service.Users.Messages.List("me").Q(query).Do()
}

// SetPrintedEmail sets an email as printed
func SetPrintedEmail(labelId string, mailId string, service *gmail.Service) (*gmail.Message, error) {
	data := &gmail.ModifyMessageRequest{AddLabelIds: []string{labelId}}
	return service.Users.Messages.Modify("me", mailId, data).Do()
}
