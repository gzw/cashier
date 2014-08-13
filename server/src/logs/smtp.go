//cashier

//@description cashier system
//@link https://github.com/gzw/cashier
//@authors gzw

package logs

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

const (
	subjectPhrase = "Diagnostic message from server"
)

type SmtpWriter struct {
	Username           string   `json:"Username"`
	Password           string   `json:"password"`
	Host               string   `json:"Host"`
	Subject            string   `json:"subject"`
	RecipientAddresses []string `json:"sendTos"`
	Level              int      `json:"level"`
}

func NewSmtpWriter() LoggerInterface {
	return &SmtpWriter{Level: LevelTrace}
}

// init SmtpWriter with json config
// config like :
// {
//	  "UserName":"example@gmail.com"
//	  "password":"password"
//    "host":"smtp.gmail.com:465"
//	  "subject":"email title"
//    "sendTos": ["email1", "email2"],
//	  "level" :LevelError
// }

func (s *SmtpWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), s)
	if err != nil {
		return err
	}
	return nil
}

// write message in smtp writer
// it will send an email with subject and this message.
func (s *SmtpWriter) WriteMsg(msg string, level int) error {
	if level < s.Level {
		return nil
	}

	hp := strings.Split(s.Host, ":")

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		s.Username,
		s.Password,
		hp[0])
	// Connect to the server, authenticate, set the sender and recipient
	// and Send the mail all in one step
	content_type := "Content-type: text/plain" + "; charset=UTF-8"
	mailmsg := []byte("To: " + strings.Join(s.RecipientAddresses, ";") + "\r\nFrom: " + s.Username + "<" + s.Username +
		">\r\nsubject: " + s.Subject + "\r\n" + content_type + "\r\n\r\n" + fmt.Sprintf(".%s", time.Now().Format("2006-01-02 15:04:05")) + msg)
	err := smtp.SendMail(s.Host, auth, s.Username, s.RecipientAddresses, mailmsg)
	return err
}

// implementing method .empty
func (s *SmtpWriter) Flush() {
	return
}

// implementing method. empty
func (s *SmtpWriter) Destroy() {
	return
}

func Init() {
	Register("smtp", NewSmtpWriter)

}
