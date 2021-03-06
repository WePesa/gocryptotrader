package main

import (
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	SMSGLOBAL_API_URL     = "http://www.smsglobal.com/http-api.php"
	ErrSMSContactNotFound = "SMS Contact not found."
	ErrSMSNotSent         = "SMS message not sent."
)

func GetEnabledSMSContacts() int {
	counter := 0
	for _, contact := range bot.config.SMS.Contacts {
		if contact.Enabled {
			counter++
		}
	}
	return counter
}

func SMSSendToAll(message string) {
	for _, contact := range bot.config.SMS.Contacts {
		if contact.Enabled {
			err := SMSNotify(contact.Number, message)
			if err != nil {
				log.Printf("Unable to send SMS to %s.\n", contact.Name)
			}
		}
	}
}

func SMSGetNumberByName(name string) string {
	for _, contact := range bot.config.SMS.Contacts {
		if contact.Name == name {
			return contact.Number
		}
	}
	return ErrSMSContactNotFound
}

func SMSNotify(to, message string) error {
	values := url.Values{}
	values.Set("action", "sendsms")
	values.Set("user", bot.config.SMS.Username)
	values.Set("password", bot.config.SMS.Password)
	values.Set("from", bot.config.Name)
	values.Set("to", to)
	values.Set("text", message)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resp, err := SendHTTPRequest("POST", SMSGLOBAL_API_URL, headers, strings.NewReader(values.Encode()))

	if err != nil {
		return err
	}

	if !StringContains(resp, "OK: 0; Sent queued message") {
		return errors.New(ErrSMSNotSent)
	}
	return nil
}
