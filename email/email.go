package email

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// EmailService wraps a gmail.Service instance.
type EmailService struct {
	*gmail.Service
}

// NewEmailService creates an email service based on the oauth2 credentials.
func NewEmailService(ctx context.Context, credentials *oauth2.Config, token *oauth2.Token) (*EmailService, error) {
	gmailService, err := gmail.NewService(ctx, option.WithTokenSource(credentials.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &EmailService{gmailService}, nil
}

// FindEmails finds all emails that match the query.
func (svc *EmailService) FindEmails(query string) ([]*gmail.Message, error) {
	res, err := svc.Users.Messages.List("me").Q(query).Do()
	return res.Messages, err
}

// GetEmailById gets the email by id.
func (svc *EmailService) GetEmailById(emailId string) (*gmail.Message, error) {
	return svc.Users.Messages.Get("me", emailId).Do()
}

// GetAttachmentBody gets the body of an attachment.
func (svc *EmailService) GetAttachmentById(emailId string, attachmentId string) (*gmail.MessagePartBody, error) {
	return svc.Users.Messages.Attachments.Get("me", emailId, attachmentId).Do()
}

// GetLabelByName gets a label by his name.
func (svc *EmailService) GetLabelByName(name string) (*gmail.Label, error) {
	// List all labels.
	res, err := svc.Users.Labels.List("me").Do()
	if err != nil {
		return nil, err
	}

	// Find label by name.
	for _, label := range res.Labels {
		if strings.EqualFold(label.Name, name) {
			return label, nil
		}
	}

	return nil, errors.New("no labels match the desired criteria")
}

// CreateLabel creates a label.
func (svc *EmailService) CreateLabel(name string) (*gmail.Label, error) {
	return svc.Users.Labels.Create("me", &gmail.Label{Name: name}).Do()
}

// SetEmailLabel sets a label for an email.
func (svc *EmailService) SetEmailLabel(emailId string, labelId string) (*gmail.Message, error) {
	return svc.Users.Messages.Modify("me", emailId, &gmail.ModifyMessageRequest{AddLabelIds: []string{labelId}}).Do()
}

// SaveAttachment saves an attachment with the same filename as sent in the email.
func (svc *EmailService) SaveAttachment(email *gmail.Message, attachment *gmail.MessagePart) error {
	attachmentBody, err := svc.Users.Messages.Attachments.Get("me", email.Id, attachment.Body.AttachmentId).Do()
	if err != nil {
		return err
	}

	// Allocate a byte array with the length of the attachment file.
	// Attachments are saved as url-encoded base64.
	attachmentFile := make([]byte, base64.URLEncoding.DecodedLen(len(attachmentBody.Data)))

	// Decode attachment in base64.
	_, err = base64.URLEncoding.Decode(attachmentFile, []byte(attachmentBody.Data))
	if err != nil {
		return err
	}

	// Write attachment with the same filename.
	return ioutil.WriteFile(attachment.Filename, attachmentFile, 0666)
}
