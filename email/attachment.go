package email

import (
	"github.com/MarcoTomasRodriguez/auto-printer/fs"
	"google.golang.org/api/gmail/v1"
	"os/exec"
)

// Gets all the attachments in a mail
func GetAttachments(mail *gmail.Message, service *gmail.Service) ([]*gmail.MessagePart, error) {
	var fullEmail *gmail.Message
	var err error

	fullEmail, err = service.Users.Messages.Get("me", mail.Id).Do()
	if err != nil {
		return nil, err
	}

	return fullEmail.Payload.Parts, nil
}

// Gets the body of an attachment
func GetAttachmentBody(attachmentId string, messageId string, service *gmail.Service) (*gmail.MessagePartBody, error) {
	return service.Users.Messages.Attachments.Get("me", messageId, attachmentId).Do()
}

// Saves the attachment to a file
func PrintAttachment(attachment *gmail.MessagePart, attachmentBody *gmail.MessagePartBody) error {
	var err error

	filename := attachment.Filename
	data := attachmentBody.Data

	if err = fs.SaveBase64File(filename, data); err != nil {
		return err
	}

	// Prints the attachment
	if err = exec.Command("lp", filename).Run(); err != nil {
		return err
	}

	if err = fs.RemoveFile(attachment.Filename); err != nil {
		return err
	}

	return nil
}
