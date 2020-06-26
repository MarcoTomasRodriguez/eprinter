package email

import (
	"fmt"
	"github.com/MarcoTomasRodriguez/auto-printer/conf"
	"google.golang.org/api/gmail/v1"
)

// Gets (if it doesn't exist, creates) the printed label, and returns his id
func GetPrintedLabelId(service *gmail.Service) (string, error) {
	// Gets the configuration and validates the used fields
	config := conf.LoadConfig()
	if err := config.ValidatePrintedLabelName(); err != nil {
		return "", err
	}

	printedLabelName := config.PrintedLabelName

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

// Gets all the unlabelled emails which the program has to print
func GetUnlabelledEmails(service *gmail.Service) (*gmail.ListMessagesResponse, error) {
	var err error
	var allowedEmails string
	var allowedEmailSubjects string

	// Loads the configuration and validates the used fields
	config := conf.LoadConfig()
	if err = config.ValidateAllowedEmails(); err != nil {
		return nil, err
	}
	if err = config.ValidateAllowedEmailSubjects(); err != nil {
		return nil, err
	}

	// Parses the allowed emails into a valid gmail query
	for _, allowedEmail := range config.AllowedEmails {
		allowedEmails = fmt.Sprintf("%s from:%s", allowedEmails, allowedEmail)
	}

	// Parses the allowed subjects into a valid gmail query
	for _, allowedEmailSubject := range config.AllowedEmailSubjects {
		allowedEmailSubjects = fmt.Sprintf("%s subject:%s", allowedEmailSubjects, allowedEmailSubject)
	}

	// Parses the query and executes it
	query := fmt.Sprintf("{%s} {%s} has:attachment has:nouserlabels", allowedEmails, allowedEmailSubjects)
	return service.Users.Messages.List("me").Q(query).Do()
}

// Sets an email as printed
func SetPrintedEmail(labelId string, mailId string, service *gmail.Service) (*gmail.Message, error) {
	data := &gmail.ModifyMessageRequest{AddLabelIds: []string{labelId}}
	return service.Users.Messages.Modify("me", mailId, data).Do()
}
