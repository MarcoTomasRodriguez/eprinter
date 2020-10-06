package email

import (
	"encoding/base64"
	"google.golang.org/api/gmail/v1"
	"os"
	"os/exec"
)

// GetAttachments gets all the attachments in an email
func GetAttachments(mail *gmail.Message, service *gmail.Service) ([]*gmail.MessagePart, error) {
	fullEmail, err := service.Users.Messages.Get("me", mail.Id).Do()
	if err != nil {
		return nil, err
	}

	return fullEmail.Payload.Parts, nil
}

// GetAttachmentBody gets the body of an attachment
func GetAttachmentBody(attachmentId string, messageId string, service *gmail.Service) (*gmail.MessagePartBody, error) {
	return service.Users.Messages.Attachments.Get("me", messageId, attachmentId).Do()
}

// PrintAttachment saves the attachment to a file
func PrintAttachment(attachment *gmail.MessagePart, attachmentBody *gmail.MessagePartBody) error {
	// Saves the attachment
	if err := saveBase64File(attachment.Filename, attachmentBody.Data); err != nil {
		return err
	}

	// Prints the attachment
	if err := exec.Command("lp", attachment.Filename).Run(); err != nil {
		return err
	}

	// Removes the attachment
	if err := os.Remove(attachment.Filename); err != nil {
		return err
	}

	return nil
}

// saveBase64File decodes the base64 data and writes it to a file
func saveBase64File(filename string, data string) error {
	// Decodes the file
	dec, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	// Creates the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writes the file
	if _, err := file.Write(dec); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}
