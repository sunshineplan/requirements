package main

import "github.com/sunshineplan/utils/mail"

var to mail.Receipts

func sendMail(subject, body string, attachments []*mail.Attachment) {
	msg := &mail.Message{
		Subject:     subject,
		Body:        body,
		Attachments: attachments,
	}
	for _, to := range to {
		msg.To = mail.Receipts{to}
		if err := dialer.Send(msg); err != nil {
			svc.Print(err)
		}
	}
}
