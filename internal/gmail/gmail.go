package gmail

import (
	"bytes"
	"mime/multipart"
	"net/smtp"
	"strings"
)

// 以下 variable 可參考 Gmail 的 smtp 設定說明
var (
	host     = "smtp.gmail.com:587"
	username = "example@gmail.com"
	password = "password"
)

// SendRawEmail 寄送信件
func SendRawEmail(to []string, body []byte) error {
	// SMTP設定
	auth := smtp.PlainAuth(host, username, password, "smtp.gmail.com")
	return smtp.SendMail(
		host,
		auth,
		username,
		to,
		body,
	)
}

// SendEmail 寄送格式化的信件
func SendEmail(to []string, subject, from string, content ...[]byte) error {
	// SMTP設定
	boundary := "hhhrrr"

	// Like Mail Header
	mailBody := &bytes.Buffer{}
	mailBody.WriteString("Subject: " + subject + "\r\n")
	mailBody.WriteString("From: " + from + "\r\n")
	mailBody.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	mailBody.WriteString(`Content-Type: multipart/mixed; boundary="` + boundary + `"` + "\r\n")

	// Mail Body Writer
	w := multipart.NewWriter(mailBody)
	defer w.Close()
	err := w.SetBoundary(boundary)
	if err != nil {
		return err
	}

	// Mail Body
	bw, err := w.CreatePart(map[string][]string{})
	if err != nil {
		return err
	}

	for i := range content {
		_, err = bw.Write(content[i])
		if err != nil {
			return err
		}

		_, err = bw.Write([]byte("\r\n"))
		if err != nil {
			return err
		}
	}
	err = SendRawEmail(to, mailBody.Bytes())
	return err
}
