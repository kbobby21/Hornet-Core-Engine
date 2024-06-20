package mail

import "net/smtp"

type ZohoMail struct {
	from     string
	password string
}

func NewZohoMail(from, password string) *ZohoMail {
	return &ZohoMail{
		from:     from,
		password: password,
	}
}

func (z *ZohoMail) SendMail(to, subject, body string) error {
	auth := smtp.PlainAuth("", z.from, z.password, "smtp.zoho.com")

	err := smtp.SendMail("smtp.zoho.com:587", auth, z.from, []string{to}, []byte(
		"Subject: "+subject+"\r\n"+
			"\r\n"+
			body,
	))

	return err
}
