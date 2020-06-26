package email

import "google.golang.org/api/gmail/v1"

// Gets the email by its id
func GetEmailById(mailId string, service *gmail.Service) (*gmail.Message, error) {
	return service.Users.Messages.Get("me", mailId).Do()
}
